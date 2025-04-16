package internal

import (
	"context"
	"fmt"
	"io"
	"strings"

	inventory "github.com/aristanetworks/cloudvision-go/api/arista/inventory.v1"
	tag "github.com/aristanetworks/cloudvision-go/api/arista/tag.v2"
	workspace "github.com/aristanetworks/cloudvision-go/api/arista/workspace.v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

type DeviceInfo struct {
	DeviceID        string
	Hostname        string
	Model           string
	Version         string
	SystemMac       string
	StreamingStatus string
	DanzEnabled     bool
	MlagEnabled     bool
}

func CreateWorkspace(ctx context.Context, conn *grpc.ClientConn, workspaceID, requestID, displayName string) {
	client := workspace.NewWorkspaceConfigServiceClient(conn)
	jsonPayload := fmt.Sprintf(`{
		"value": {
			"displayName": "%s",
			"key": {"workspaceId": "%s"},
			"requestParams": {"requestId": "%s"}
		}
	}`, displayName, workspaceID, requestID)

	var req workspace.WorkspaceConfigSetRequest
	if err := protojson.Unmarshal([]byte(jsonPayload), &req); err != nil {
		panic(fmt.Sprintf("❌ Erreur parsing JSON workspace : %v", err))
	}

	resp, err := client.Set(ctx, &req)
	if err != nil {
		panic(fmt.Sprintf("❌ Erreur création workspace : %v", err))
	}
	fmt.Printf("✅ Workspace créé : %s\n", protojson.Format(resp))
}

func CreateTag(ctx context.Context, conn *grpc.ClientConn, workspaceID, label, value string, elementType, elementSubType int) {
	client := tag.NewTagConfigServiceClient(conn)
	jsonPayload := fmt.Sprintf(`{
		"value":{
			"remove":false,
			"key":{
				"workspaceId":"%s",
				"elementType":%d,
				"label":"%s",
				"value":"%s",
				"elementSubType":%d
			}
		}
	}`, workspaceID, elementType, label, value, elementSubType)

	var req tag.TagConfigSetRequest
	if err := protojson.Unmarshal([]byte(jsonPayload), &req); err != nil {
		panic(fmt.Sprintf("❌ Erreur parsing JSON tag : %v", err))
	}

	resp, err := client.Set(ctx, &req)
	if err != nil {
		panic(fmt.Sprintf("❌ Erreur ajout tag : %v", err))
	}
	fmt.Printf("🏷️  Tag ajouté : %s\n", protojson.Format(resp))
}

func AssignTagToDevice(ctx context.Context, conn *grpc.ClientConn, workspaceID, deviceID, label, value string, elementType, elementSubType int) {
	client := tag.NewTagAssignmentConfigServiceClient(conn)
	jsonPayload := fmt.Sprintf(`{
		"value": {
			"remove": false,
			"key": {
				"workspaceId": "%s",
				"elementType": %d,
				"elementSubType": %d,
				"label": "%s",
				"value": "%s",
				"deviceId": "%s"
			}
		}
	}`, workspaceID, elementType, elementSubType, label, value, deviceID)

	var req tag.TagAssignmentConfigSetRequest
	if err := protojson.Unmarshal([]byte(jsonPayload), &req); err != nil {
		panic(fmt.Sprintf("❌ Erreur parsing tag assignment JSON : %v", err))
	}

	_, err := client.Set(ctx, &req)
	if err != nil {
		panic(fmt.Sprintf("❌ Erreur assignation tag au device : %v", err))
	}

	fmt.Printf("📌 Tag '%s=%s' assigné à device %s\n", label, value, deviceID)
}

func ReadInventory(ctx context.Context, conn *grpc.ClientConn, model string) []DeviceInfo {
	client := inventory.NewDeviceServiceClient(conn)
	var req inventory.DeviceStreamRequest

	// ➤ Si un filtre model est fourni, on le passe dans la requête
	if model != "" {
		fmt.Println("passage dans model")
		data := fmt.Sprintf(`{"partialEqFilter":[{"modelName":"%s"}]}`, model)
		if err := protojson.Unmarshal([]byte(data), &req); err != nil {
			panic(fmt.Sprintf("Erreur parsing JSON filtre modèle : %v", err))
		}
	}
	stream, err := client.GetAll(ctx, &req)

	if err != nil {
		panic(fmt.Sprintf("❌ Erreur stream inventaire : %v", err))
	}

	var devices []DeviceInfo
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(fmt.Sprintf("❌ Erreur lecture stream : %v", err))
		}

		val := res.GetValue()
		features := val.GetExtendedAttributes().GetFeatureEnabled()
		devices = append(devices, DeviceInfo{
			DeviceID:        val.GetKey().GetDeviceId().GetValue(),
			Hostname:        val.GetHostname().GetValue(),
			Model:           val.GetModelName().GetValue(),
			Version:         val.GetSoftwareVersion().GetValue(),
			SystemMac:       val.GetSystemMacAddress().GetValue(),
			StreamingStatus: val.GetStreamingStatus().String(),
			DanzEnabled:     features["Danz"],
			MlagEnabled:     features["Mlag"],
		})
	}
	return devices
}

func GetWorkspacesByState(ctx context.Context, conn *grpc.ClientConn, stateName string) []struct {
	ID          string
	DisplayName string
	State       string
} {
	stateMap := map[string]int{
		"UNSPECIFIED": 0,
		"UNRECOGNIZED": -1,
		"PENDING":     1,
		"SUBMITTED":   2,
		"ABANDONED":   3,
		"CONFLICTS":   4,
		"ROLLED_BACK": 5,
	}

	var req workspace.WorkspaceStreamRequest

	// Appliquer le filtre uniquement si on a une valeur explicite
	if stateName != "" && strings.ToUpper(stateName) != "NONE" {
		stateValue, ok := stateMap[strings.ToUpper(stateName)]
		if !ok {
			panic(fmt.Sprintf("❌ État invalide : %s", stateName))
		}

		// Construire le JSON de filtre dynamiquement
		data := fmt.Sprintf(`{"partialEqFilter":[{"state":%d}]}`, stateValue)
		if err := protojson.Unmarshal([]byte(data), &req); err != nil {
			panic(fmt.Sprintf("Erreur parsing JSON request : %v", err))
		}
	}

	client := workspace.NewWorkspaceServiceClient(conn)
	stream, err := client.GetAll(ctx, &req)
	if err != nil {
		panic(fmt.Sprintf("Erreur stream : %v", err))
	}

	var results []struct {
		ID          string
		DisplayName string
		State       string
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(fmt.Sprintf("Erreur lecture : %v", err))
		}
		val := res.GetValue()
		results = append(results, struct {
			ID          string
			DisplayName string
			State       string
		}{
			ID:          val.GetKey().GetWorkspaceId().GetValue(),
			DisplayName: val.GetDisplayName().GetValue(),
			State:       val.GetState().String(),
		})
	}
	return results
}
