package spotinst

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/spotinst/spotinst-sdk-go/service/elastigroup/providers/azure"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/commons"
	"log"
	"testing"
)

func createElastigroupAzureResourceName(name string) string {
	return fmt.Sprintf("%v.%v", string(commons.ElastigroupAzureResourceName), name)
}

func testElastigroupAzureDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != string(commons.ElastigroupAzureResourceName) {
			continue
		}
		input := &azure.ReadGroupInput{GroupID: spotinst.String(rs.Primary.ID)}
		resp, err := client.elastigroup.CloudProviderAzure().Read(context.Background(), input)
		if err == nil && resp != nil && resp.Group != nil {
			return fmt.Errorf("group still exists")
		}
	}
	return nil
}

func testCheckElastigroupAzureAttributes(group *azure.Group, expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if spotinst.StringValue(group.Name) != expectedName {
			return fmt.Errorf("bad content: %v", group.Name)
		}
		return nil
	}
}

func testCheckElastigroupAzureExists(group *azure.Group, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource ID is set")
		}
		client := testAccProvider.Meta().(*Client)
		input := &azure.ReadGroupInput{GroupID: spotinst.String(rs.Primary.ID)}
		resp, err := client.elastigroup.CloudProviderAzure().Read(context.Background(), input)
		if err != nil {
			return err
		}
		if spotinst.StringValue(resp.Group.Name) != rs.Primary.Attributes["name"] {
			return fmt.Errorf("Group not found: %+v,\n %+v\n", resp.Group, rs.Primary.Attributes)
		}
		*group = *resp.Group
		return nil
	}
}

type AzureGroupConfigMetadata struct {
	variables            string
	groupName            string
	vmSizes              string
	launchConfig         string
	strategy             string
	image                string
	loadBalancers        string
	network              string
	login                string
	fieldsToAppend       string
	updateBaselineFields bool
}

func createElastigroupAzureTerraform(gcm *AzureGroupConfigMetadata) string {
	if gcm == nil {
		return ""
	}

	if gcm.vmSizes == "" {
		gcm.vmSizes = testAzureVMSizesGroupConfig_Create
	}

	if gcm.launchConfig == "" {
		gcm.launchConfig = testAzureLaunchConfigurationGroupConfig_Create
	}

	if gcm.strategy == "" {
		gcm.strategy = testAzureStrategyGroupConfig_Create
	}

	if gcm.image == "" {
		gcm.image = testAzureImageGroupConfig_Create
	}

	if gcm.loadBalancers == "" {
		gcm.loadBalancers = testAzureLoadBalancersGroupConfig_Create
	}

	if gcm.network == "" {
		gcm.network = testAzureNetworkGroupConfig_Create
	}

	if gcm.login == "" {
		gcm.login = testAzureLoginGroupConfig_Create
	}

	template := ""
	if gcm.updateBaselineFields {
		format := testBaselineAzureGroupConfig_Update
		template = fmt.Sprintf(format,
			gcm.groupName,
			gcm.groupName,
			gcm.vmSizes,
			gcm.launchConfig,
			gcm.strategy,
			gcm.image,
			gcm.loadBalancers,
			gcm.network,
			gcm.login,
			gcm.fieldsToAppend,
		)
	} else {
		format := testBaselineAzureGroupConfig_Create
		template = fmt.Sprintf(format,
			gcm.groupName,
			gcm.groupName,
			gcm.vmSizes,
			gcm.launchConfig,
			gcm.strategy,
			gcm.fieldsToAppend,
		)
	}

	if gcm.variables != "" {
		template = gcm.variables + "\n" + template
	}

	log.Printf("Terraform [%v] template:\n%v", gcm.groupName, template)
	return template
}

// region Elastigroup Azure: Baseline
func TestAccSpotinstElastigroupAzure_Baseline(t *testing.T) {
	groupName := "eg-baseline"
	resourceName := createElastigroupAzureResourceName(groupName)

	var group azure.Group
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testElastigroupAzureDestroy,

		Steps: []resource.TestStep{
			{
				Config: createElastigroupAzureTerraform(&AzureGroupConfigMetadata{groupName: groupName}),
				Check: resource.ComposeTestCheckFunc(
					testCheckElastigroupAzureExists(&group, resourceName),
					testCheckElastigroupAzureAttributes(&group, groupName),
					resource.TestCheckResourceAttr(resourceName, "max_size", "0"),
					resource.TestCheckResourceAttr(resourceName, "min_size", "0"),
					resource.TestCheckResourceAttr(resourceName, "desired_capacity", "0"),
				),
			},
			{
				Config: createElastigroupAzureTerraform(&AzureGroupConfigMetadata{groupName: groupName, updateBaselineFields: true}),
				Check: resource.ComposeTestCheckFunc(
					testCheckElastigroupAzureExists(&group, resourceName),
					testCheckElastigroupAzureAttributes(&group, groupName),
					resource.TestCheckResourceAttr(resourceName, "max_size", "0"),
					resource.TestCheckResourceAttr(resourceName, "min_size", "0"),
					resource.TestCheckResourceAttr(resourceName, "desired_capacity", "0"),
				),
			},
		},
	})
}

const testBaselineAzureGroupConfig_Create = `
resource "` + string(commons.ElastigroupAzureResourceName) + `" "%v" {

 name 				 = "%v"
 product 			 = "Linux"
 resource_group_name = "alex-test"

 // --- CAPACITY ------------
 max_size 		  = 0
 min_size 		  = 0
 desired_capacity = 0
 // -------------------------
 
 %v
 %v
 %v
 %v
}

`

const testBaselineAzureGroupConfig_Update = `
resource "` + string(commons.ElastigroupAzureResourceName) + `" "%v" {

 name 				 = "%v"
 product 			 = "Linux"
 resource_group_name = "alex-test"

 // --- CAPACITY ------------
 max_size 		  = 0
 min_size 		  = 0
 desired_capacity = 0
 // -------------------------
 
 %v
 %v
 %v
 %v
}

`

// endregion

// region Azure Elastigroup: VM Sizes

func TestAccSpotinstElastigroupAzure_VMSizes(t *testing.T) {
	groupName := "eg-azure-vm-sizes"
	resourceName := createElastigroupAzureResourceName(groupName)

	var group azure.Group
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testElastigroupAzureDestroy,

		Steps: []resource.TestStep{
			{
				Config: createElastigroupAzureTerraform(&AzureGroupConfigMetadata{
					groupName: groupName,
					vmSizes:   testAzureVMSizesGroupConfig_Create,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckElastigroupAzureExists(&group, resourceName),
					testCheckElastigroupAzureAttributes(&group, groupName),
					resource.TestCheckResourceAttr(resourceName, "od_sizes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "od_sizes.0", "basic_a1"),
					resource.TestCheckResourceAttr(resourceName, "low_priority_sizes.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "low_priority_sizes.0", "basic_a2"),
					resource.TestCheckResourceAttr(resourceName, "low_priority_sizes.1", "basic_a1"),
					resource.TestCheckResourceAttr(resourceName, "low_priority_sizes.2", "basic_a3"),
				),
			},
			{
				Config: createElastigroupAzureTerraform(&AzureGroupConfigMetadata{
					groupName: groupName,
					vmSizes:   testAzureVMSizesGroupConfig_Update,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckElastigroupAzureExists(&group, resourceName),
					testCheckElastigroupAzureAttributes(&group, groupName),
					resource.TestCheckResourceAttr(resourceName, "od_sizes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "od_sizes.0", "basic_a2"),
					resource.TestCheckResourceAttr(resourceName, "low_priority_sizes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "low_priority_sizes.0", "basic_a2"),
				),
			},
		},
	})
}

const testAzureVMSizesGroupConfig_Create = `
// --- VM SIZES --------------------------------
 od_sizes           = ["basic_a1"]
 low_priority_sizes = ["basic_a2", "basic_a1", "basic_a3"]
// ---------------------------------------------
`

const testAzureVMSizesGroupConfig_Update = `
// --- VM SIZES --------------------------------
 od_sizes           = ["basic_a2"]
 low_priority_sizes = ["basic_a2"]
// ---------------------------------------------
`

// endregion

// region Azure Elastigroup: Launch Configuration

func TestAccSpotinstElastigroupAzure_LaunchConfiguration(t *testing.T) {
	groupName := "eg-azure-launch-configuration"
	resourceName := createElastigroupAzureResourceName(groupName)

	var group azure.Group
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testElastigroupAzureDestroy,

		Steps: []resource.TestStep{
			{
				Config: createElastigroupAzureTerraform(&AzureGroupConfigMetadata{
					groupName: groupName,
					vmSizes:   testAzureLaunchConfigurationGroupConfig_Create,
				}),
				Check: resource.ComposeTestCheckFunc(),
			},
			{
				Config: createElastigroupAzureTerraform(&AzureGroupConfigMetadata{
					groupName: groupName,
					vmSizes:   testAzureLaunchConfigurationGroupConfig_Update,
				}),
				Check: resource.ComposeTestCheckFunc(),
			},
		},
	})
}

const testAzureLaunchConfigurationGroupConfig_Create = `
// --- LAUNCH CONFIGURATION --------------------
// user_data = ""

// ---------------------------------------------
`

const testAzureLaunchConfigurationGroupConfig_Update = `
// --- LAUNCH CONFIGURATION --------------------
// user_data = ""

// ---------------------------------------------
`

// endregion

// region Azure Elastigroup: Strategy

//func TestAccSpotinstElastigroupAzure_Strategy(t *testing.T) {
//	groupName := "eg-azure-strategy"
//	resourceName := createElastigroupAzureResourceName(groupName)
//
//	var group azure.Group
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() { testAccPreCheck(t) },
//		Providers: TestAccProviders,
//		CheckDestroy: testElastigroupAzureDestroy,
//
//		Steps: []resource.TestStep{
//			{
//				Config: createElastigroupAzureTerraform(&AzureGroupConfigMetadata{
//					groupName: groupName,
//					strategy: testAzureStrategyGroupConfig_Create,
//				}),
//				Check: resource.ComposeTestCheckFunc(),
//			},
//			{
//				Config: createElastigroupAzureTerraform(&AzureGroupConfigMetadata{
//					groupName: groupName,
//					strategy: testAzureStrategyGroupConfig_Update,
//				}),
//				Check: resource.ComposeTestCheckFunc(),
//			},
//			{
//				Config: createElastigroupAzureTerraform(&AzureGroupConfigMetadata{
//					groupName: groupName,
//					strategy: testAzureStrategyGroupConfig_EmptyFields,
//				}),
//				Check: resource.ComposeTestCheckFunc(),
//			},
//		},
//	})
//}

const testAzureStrategyGroupConfig_Create = `
// --- STRATEGY --------------------------------
 low_priority_percentage = 50
 on_demand_count         = 1
 draining_timeout        = 180
// ---------------------------------------------
`

const testAzureStrategyGroupConfig_Update = `
// --- STRATEGY --------------------------------
 low_priority_percentage = 50
 on_demand_count         = 1
 draining_timeout        = 180
// ---------------------------------------------
`

const testAzureStrategyGroupConfig_EmptyFields = `
// --- STRATEGY --------------------------------
// ---------------------------------------------
`

// endregion

// region Azure Elastigroup: Image

//func TestAccSpotinstElastigroupAzure_Image(t *testing.T) {
//	groupName := "eg-azure-strategy"
//	resourceName := createElastigroupAzureResourceName(groupName)
//
//	var group azure.Group
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() { testAccPreCheck(t) },
//		Providers: TestAccProviders,
//		CheckDestroy: testElastigroupAzureDestroy,
//
//		Steps: []resource.TestStep{
//			{
//				Config: createElastigroupAzureTerraform(&AzureGroupConfigMetadata{
//					groupName: groupName,
//					strategy: testAzureStrategyGroupConfig_Create,
//				}),
//				Check: resource.ComposeTestCheckFunc(),
//			},
//			{
//				Config: createElastigroupAzureTerraform(&AzureGroupConfigMetadata{
//					groupName: groupName,
//					strategy: testAzureImageGroupConfig_Update,
//				}),
//				Check: resource.ComposeTestCheckFunc(),
//			},
//			{
//				Config: createElastigroupAzureTerraform(&AzureGroupConfigMetadata{
//					groupName: groupName,
//					strategy: testAzureImageGroupConfig_EmptyFields,
//				}),
//				Check: resource.ComposeTestCheckFunc(),
//			},
//		},
//	})
//}

const testAzureImageGroupConfig_Create = `
// --- IMAGES --------------------------------
  image = {
    marketplace = {
      publisher = "Canonical"
      offer = "UbuntuServer"
      sku = "16.04-LTS"
    }
  }
// ---------------------------------------------
`

const testAzureImageGroupConfig_Update = `
// --- IMAGES --------------------------------
  image = {
    marketplace = {
      publisher = ""
      offer = ""
      sku = ""
    }
  }
// ---------------------------------------------
`

const testAzureImageGroupConfig_EmptyFields = `
// --- IMAGES --------------------------------
// ---------------------------------------------
`

// endregion

// region Azure Elastigroup: Load Balancers

//func TestAccSpotinstElastigroupAzure_Image(t *testing.T) {
//	groupName := "eg-azure-strategy"
//	resourceName := createElastigroupAzureResourceName(groupName)
//
//	var group azure.Group
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() { testAccPreCheck(t) },
//		Providers: TestAccProviders,
//		CheckDestroy: testElastigroupAzureDestroy,
//
//		Steps: []resource.TestStep{
//			{
//				Config: createElastigroupAzureTerraform(&AzureGroupConfigMetadata{
//					groupName: groupName,
//					strategy: testAzureStrategyGroupConfig_Create,
//				}),
//				Check: resource.ComposeTestCheckFunc(),
//			},
//			{
//				Config: createElastigroupAzureTerraform(&AzureGroupConfigMetadata{
//					groupName: groupName,
//					strategy: testAzureImageGroupConfig_Update,
//				}),
//				Check: resource.ComposeTestCheckFunc(),
//			},
//			{
//				Config: createElastigroupAzureTerraform(&AzureGroupConfigMetadata{
//					groupName: groupName,
//					strategy: testAzureImageGroupConfig_EmptyFields,
//				}),
//				Check: resource.ComposeTestCheckFunc(),
//			},
//		},
//	})
//}

const testAzureLoadBalancersGroupConfig_Create = `
// --- LOAD BALANCERS --------------------------
  load_balancers = [
    {
      type = "MULTAI_TARGET_SET"
      balancer_id = "lb-0be85d6aa269"
      target_set_id = "ts-ae9c9603c365"
      auto_weight = true
    }
  ]
// ---------------------------------------------
`

const testAzureLoadBalancersGroupConfig_Update = `
// --- LOAD BALANCERS --------------------------
  load_balancers = [
    {
      type = "MULTAI_TARGET_SET"
      balancer_id = "lb-0be85d6aa269"
      target_set_id = "ts-ae9c9603c365"
      auto_weight = false
    }
  ]
// ---------------------------------------------
`

const testAzureLoadBalancersGroupConfig_EmptyFields = `
// --- LOAD BALANCERS --------------------------
// ---------------------------------------------
`

// endregion

// region Azure Elastigroup: Network

//func TestAccSpotinstElastigroupAzure_Network(t *testing.T) {
//	groupName := "eg-azure-strategy"
//	resourceName := createElastigroupAzureResourceName(groupName)
//
//	var group azure.Group
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() { testAccPreCheck(t) },
//		Providers: TestAccProviders,
//		CheckDestroy: testElastigroupAzureDestroy,
//
//		Steps: []resource.TestStep{
//			{
//				Config: createElastigroupAzureTerraform(&AzureGroupConfigMetadata{
//					groupName: groupName,
//					strategy: testAzureStrategyGroupConfig_Create,
//				}),
//				Check: resource.ComposeTestCheckFunc(),
//			},
//			{
//				Config: createElastigroupAzureTerraform(&AzureGroupConfigMetadata{
//					groupName: groupName,
//					strategy: testAzureImageGroupConfig_Update,
//				}),
//				Check: resource.ComposeTestCheckFunc(),
//			},
//			{
//				Config: createElastigroupAzureTerraform(&AzureGroupConfigMetadata{
//					groupName: groupName,
//					strategy: testAzureImageGroupConfig_EmptyFields,
//				}),
//				Check: resource.ComposeTestCheckFunc(),
//			},
//		},
//	})
//}

const testAzureNetworkGroupConfig_Create = `
// --- NETWORK ---------------------------------
  network = {
    virtual_network_name = "alex-test-netwrk"
    subnet_name = "alex-test-subnet"                 
    resource_group_name = "alex-test"         
    assign_public_ip = true                
  }
// ---------------------------------------------
`

const testAzureNetworkGroupConfig_Update = `
// --- NETWORK ---------------------------------
  network = {
    virtual_network_name = "alex-test-netwrk"
    subnet_name = "alex-test-subnet"                 
    resource_group_name = "alex-test"         
    assign_public_ip = true                
  }
// ---------------------------------------------
`

const testAzureNetworkGroupConfig_EmptyFields = `
// --- NETWORK ---------------------------------
// ---------------------------------------------
`

// endregion

// region Azure Elastigroup: Network

//func TestAccSpotinstElastigroupAzure_Network(t *testing.T) {
//	groupName := "eg-azure-strategy"
//	resourceName := createElastigroupAzureResourceName(groupName)
//
//	var group azure.Group
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() { testAccPreCheck(t) },
//		Providers: TestAccProviders,
//		CheckDestroy: testElastigroupAzureDestroy,
//
//		Steps: []resource.TestStep{
//			{
//				Config: createElastigroupAzureTerraform(&AzureGroupConfigMetadata{
//					groupName: groupName,
//					strategy: testAzureStrategyGroupConfig_Create,
//				}),
//				Check: resource.ComposeTestCheckFunc(),
//			},
//			{
//				Config: createElastigroupAzureTerraform(&AzureGroupConfigMetadata{
//					groupName: groupName,
//					strategy: testAzureImageGroupConfig_Update,
//				}),
//				Check: resource.ComposeTestCheckFunc(),
//			},
//			{
//				Config: createElastigroupAzureTerraform(&AzureGroupConfigMetadata{
//					groupName: groupName,
//					strategy: testAzureImageGroupConfig_EmptyFields,
//				}),
//				Check: resource.ComposeTestCheckFunc(),
//			},
//		},
//	})
//}

const testAzureLoginGroupConfig_Create = `
// --- LOGIN ---------------------------------
  login = {
    user_name = "alex-test"
    ssh_public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDLn7RIjgivW2nWoh56XV2wpDKjjWFk92UgfTsqL8qYI0lGCJuoGeMlje1VWyAemMGZZsK5et8j3/caZsVd1Hypui3xV+tRAmtnyqVCjDGYsBQIMKoDzLrrZz7/s2WNKbMegOgQ+8YxXLhxuS5YGKhNjvxC2kJCe1HkPAPvx03kzNGmxxv6pt5TaQPXUqVxfWoeoaLRDcL8Ns2kikZC6v2cfY/PcmwoYd7XlVuILLTMNF6ujOUsX9kHt/910dEW66iZpc+PjHnKuAu/5238lssiUZULTHWbjE09MG8kHIiZq9Z9HgmAS++YLUc2G9InBqiLXMbie4S9qMcp+crl1oG/"
  }
// -------------------------------------------
`

const testAzureLoginGroupConfig_Update = `
// --- LOGIN ---------------------------------
  login = {
    user_name = "alex-test"
    ssh_public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDLn7RIjgivW2nWoh56XV2wpDKjjWFk92UgfTsqL8qYI0lGCJuoGeMlje1VWyAemMGZZsK5et8j3/caZsVd1Hypui3xV+tRAmtnyqVCjDGYsBQIMKoDzLrrZz7/s2WNKbMegOgQ+8YxXLhxuS5YGKhNjvxC2kJCe1HkPAPvx03kzNGmxxv6pt5TaQPXUqVxfWoeoaLRDcL8Ns2kikZC6v2cfY/PcmwoYd7XlVuILLTMNF6ujOUsX9kHt/910dEW66iZpc+PjHnKuAu/5238lssiUZULTHWbjE09MG8kHIiZq9Z9HgmAS++YLUc2G9InBqiLXMbie4S9qMcp+crl1oG/"
  }
// -------------------------------------------
`

const testAzureLoginGroupConfig_EmptyFields = `
// --- LOGIN ---------------------------------
// -------------------------------------------
`
