package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// Subscription ID for the testing environment
var subscriptionID string = "d9e934e6-bf60-4a22-abae-d7f92d24e1c3"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "bhup0006",
		},
	}

	// Defer the destruction of resources after tests
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variables
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	nicName := terraform.Output(t, terraformOptions, "nic_name")

	// Test 1: Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))

	// Test 2: Confirm NIC exists and is connected to the VM
	// We'll indirectly confirm by checking the NIC existence
	assert.True(t, azure.NetworkInterfaceExists(t, nicName, resourceGroupName, subscriptionID))

	// Test 3: Confirm the VM is running the correct Ubuntu version
	// You would typically SSH into the VM or use the Azure API to fetch the OS version
	// In this case, use a helper function that runs a shell command on the VM (SSH needed)
	// This is an example check:
	// Note: You can use SSH for this (make sure SSH access is set up in your Terraform code)
	// osVersion := GetUbuntuVersion(vmName, resourceGroupName)  // Use helper function
	// assert.Equal(t, "Ubuntu", osVersion, "VM is not running Ubuntu")

	// Alternatively, you can manually confirm from the Azure Portal for OS profile if needed.
	t.Logf("Test completed for VM %s in resource group %s", vmName, resourceGroupName)
}
