package client

import (
	"context"
	"fmt"

	caphv1beta1 "github.com/syself/cluster-api-provider-hetzner/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// A private IP might be assigned to a Hetzner Bare Metal server (HBMS), using the private-ip label
// in the HetznerBareMetalHost resource.
// This function is responsible for checking whether the assignment exists. And if yes, then it
// returns that assigned private IP.
func GetHBMSPrivateIP(clusterClient client.Client, serverID string) (*string, error) {
	// Try getting the corresponding HetznerBareMetalHost.
	hetznerBareMetalHost := &caphv1beta1.HetznerBareMetalHost{
		ObjectMeta: metav1.ObjectMeta{
			Name: serverID,

			// TODO : Take the capi-cluster namespace via a CLI flag.
			Namespace: "capi-cluster",
		},
	}
	err := clusterClient.Get(context.Background(),
		types.NamespacedName{
			Name:      serverID,
			Namespace: "capi-cluster",
		},
		hetznerBareMetalHost,
	)
	if err != nil {
		return nil, fmt.Errorf("failed getting HetznerBareMetalHost %s : %v", serverID, err)
	}

	// Check whether the HetznerBareMetalHost resource has the private-ip label.
	// When yes, it indicates that we have assigned a private IP to the corresponding HBMS.
	if privateIP, ok := hetznerBareMetalHost.Labels["private-ip"]; ok {
		return &privateIP, nil
	}
	return nil, nil
}
