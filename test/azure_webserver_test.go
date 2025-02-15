package test
 
import (
    "testing"
 
    "github.com/gruntwork-io/terratest/modules/azure"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)
 
// Subscription ID for Azure
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
 
    // Confirm VM exists
    assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))
}
 
func TestNetworkInterfaceExists(t *testing.T) {
    t.Parallel()
 
    terraformOptions := &terraform.Options{
        TerraformDir: "../",
    }
 
    terraform.InitAndApply(t, terraformOptions)
 
    nicID := terraform.Output(t, terraformOptions, "nic_id")
    vmID := terraform.Output(t, terraformOptions, "vm_id")
 
    // Test that the NIC exists
    nicExists := azure.NetworkInterfaceExists(t, nicID, subscriptionID)
    assert.True(t, nicExists, "Network Interface should exist")
 
    // Check if NIC is connected to the VM
    vmNicConnection := azure.CheckNICConnection(t, vmID, nicID, subscriptionID)
    assert.True(t, vmNicConnection, "NIC should be connected to the VM")
 
    terraform.Destroy(t, terraformOptions)
}
 
func TestUbuntuVersion(t *testing.T) {
    t.Parallel()
 
    terraformOptions := &terraform.Options{
        TerraformDir: "../",
    }
 
    terraform.InitAndApply(t, terraformOptions)
 
    publicIP := terraform.Output(t, terraformOptions, "public_ip")
 
    sshUser := "azureuser"
    sshKeyPath := "file("~/.ssh/id_rsa.pub")"
    versionCheckCmd := "lsb_release -a"
    output := terraform.RunSSHCommand(t, publicIP, sshUser, sshKeyPath, versionCheckCmd)
 
    assert.Contains(t, output, "Ubuntu 20.04", "VM should be running Ubuntu 20.04")
 
    terraform.Destroy(t, terraformOptions)
}