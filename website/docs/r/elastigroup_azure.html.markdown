---
layout: "spotinst"
page_title: "Spotinst: azure_elastigroup"
sidebar_current: "docs-do-resource-azure_elastigroup"
description: |-
Provides a Spotinst Azure elastigroup resource.
---

# spotinst\_azure\_elastigroup

Provides a Spotinst Azure elastigroup resource.

## Example Usage

```hcl
resource "azure_elastigroup" "test_azure_group" {
  name                = "azure_elastigroup"
  resource_group_name = "spotinst-azure"
  region              = "eastus"
  product             = "Linux"
  user_data = ""

  // --- CAPACITY -----------------------------------------------------
  min_size         = 0
  max_size         = 1
  desired_capacity = 0
  // ------------------------------------------------------------------

  // --- LAUNCH CONFIGURATION -----------------------------------
  image = {
    // market image
    publisher = "OpenLogic"
    offer     = "CentOS"
    sku       = "7.3"

    // custom image types are required when market is not specified
    # resource_group_name = "spotinst_azure"
    # image_name          = "custom-image-name"
  }
  // ------------------------------------------------------------

  // --- INSTANCE TYPES -------------------------------------------------------
  od_sizes           = ["standard_a1_v1", "standard_a1_v2"]
  low_priority_sizes = ["standard_a1_v1", "standard_a1_v2"]
  // --------------------------------------------------------------------------

  // --- STRATEGY -------------------------------------------------------------
  low_priority_percentage = 100
  od_count                = 1
  draining_timeout        = 300
  // -------------------------------------------------------------------------

  // --- SCALING POLICIES --------------------------------------------------
  scaling_up_policy = [{
     policy_name = "policy-name"
     metric_name = "CPUUtilization"
     statistic   = "average"
     threshold   = 10
     period      = 60
     cooldown    = 300
     operator    = "gte"

     evaluation_periods = 10

     action = {
         type       = "adjustment"
         adjustment = 2
     }
   }]

  scaling_down_policy = [{...}] // same as scaling_up_policy
  // -----------------------------------------------------------------------


  // --- LOAD BALANCERS ------------------------------
  load_balancers = [{
    type          = "MULTAI_TARGET_SET"
    balancer_id   = "lb-1ee2e3q"
    target_set_id = "ts-3eq"
    auto_weight   = true
  }]

  // --------------------------------------------------

  // --- HEALTH-CHECKS ------------------------------------
  health_check_type         = "INSTANCE_STATE"
  health_check_grace_period = 120
  auto_healing              = true
  // ------------------------------------------------------

  // --- NETWORK ------------------------------------
  network = {
    virtual_network_name = "vname"
    subnet_name          = "my-subnet-name"
    resource_group_name  = "subnetResourceGroup"
    assign_public_up     = true
  }
  // ------------------------------------------------

  // --- LOGIN ------------------------------------
  user_name      = "admin"
  ssh_public_key = "33a2s1f3g5a1df5g1ad3f2g1adfg56dfg=="
  password       = "****************"
  // ----------------------------------------------
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The group name.
* `region` - (Required) The region your Azure group will be created in.
* `resource_group_name` - (Required) Name of the Resource Group for Elastigroup.
* `product` - (Required) Operation system type. Valid values: `"Linux"`, `"Windows"`.
* `max_size` - (Required) The maximum number of instances the group should have at any time.
* `min_size` - (Required) The minimum number of instances the group should have at any time.
* `desired_capacity` - (Required) The desired number of instances the group should have at any time.

* `image` - (Required) Image of a VM. An image is a template for creating new VMs. Choose from Azure image catalogue (marketplace) or use a custom image.
* `publisher` - (Optional) Image publisher. Required if resource_group_name is not specified.
* `offer` - (Optional) Name of the image to use. Required if publisher is specified.
* `user_data` - (Optional) Base64-encoded MIME user data to make available to the instances.
* `sku` - (Optional) Image’s Stock Keeping Unit, which is the specific version of the image. Required if publisher is specified.
* `resource_group_name` - (Optional) Name of Resource Group for custom image. Required if publisher not specified.
* `image_name` - (Optional) Name of the custom image. Required if resource_group_name is specified.

* `od_sizes` - (Required) Available On-Demand sizes
* `low_priority_sizes` - (Required) Available Low-Priority sizes.

* `low_priority_percentage` - (Optional, Default `100`) Percentage of Low Priority instances to maintain. Required if `od_count` is not specified.
* `od_count` - (Optional) Number of On-Demand instances to maintain. Required if low_priority_percentage is not specified.
* `draining_timeout` - (Optional, Default `120`) Time (seconds) to allow the instance to be drained from incoming TCP connections and detached from MLB before terminating it during a scale-down operation.

* `scaling_up_policy` - (Optional) Contains scaling policies for scaling the Elastigroup up.
* `scaling_down_policy` - (Optional) Contains scaling policies for scaling the Elastigroup down.
* `policy_name` - (Optional) Name of scaling policy.
* `metric_name` - (Optional) Metric to monitor. Valid values: “Percentage CPU”, “Network In”, “Network Out”, “Disk Read Bytes”, “Disk Write Bytes”, “Disk Write Operations/Sec”, “Disk Read Operations/Sec”.
* `statistic` - (Optional) Statistic by which to evaluate the selected metric. Valid values: “AVERAGE”, “SAMPLE_COUNT”, “SUM”, “MINIMUM”, “MAXIMUM”, “PERCENTILE”, “COUNT”.
* `threshold` - (Optional) The value at which the scaling action is triggered.
* `period` - (Optional) Amount of time (seconds) for which the threshold must be met in order to trigger the scaling action.
* `evaluation_periods` - (Optional) Number of consecutive periods in which the threshold must be met in order to trigger a scaling action.
* `cooldown` - (Optional) Time (seconds) to wait after a scaling action before resuming monitoring.
* `operator` - (Optional) The operator used to evaluate the threshold against the current metric value. Valid values: “gt” (greater than), “get” (greater-than or equal), “lt” (less than), “lte” (less than or equal).
* `action` - (Optional) Scaling action to take when the policy is triggered.
* `type` - (Optional) Type of scaling action to take when the scaling policy is triggered. Valid values: “adjustment”, “setMinTarget”, “updateCapacity”, “percentageAdjustment”
* `adjustment` - (Optional) Value to which the action type will be adjusted. Required if using “numeric” or “percentageAdjustment” action types.
* `max_target_capacity` - (Optional) Value to adjust the capacity to. Required if using “setMinTarget” as action type.
* `target` - (Optional) Desired number of instances. Required if using ‘updateCpacity’ as action type and ‘minimum’ and ‘maximum’ are defined.
* `minimum` - (Optional) Lower limit of instances that you can scale to. Required if using ‘updateCpacity’ as action type and ‘target’ and ‘maximum’ are defined.
* `maximum` - (Optional)Upper limit of instances that you can scale to. Required if using ‘updateCpacity’ as action type and ‘minimum’ and ‘target’ are defined.

* `load_balancers` - (Required) Describes a set of one or more classic load balancer target groups and/or Multai load balancer target sets.
* `type` - (Required) The resource type. Valid values: CLASSIC, TARGET_GROUP, MULTAI_TARGET_SET.
* `balancer_id` - (Required) The balancer ID.
* `target_set_id` - (Required) The scale set ID associated with the load balancer.
* `auto_weight` - (Optional, Default: `false`)

* `health_check_type` - (Optional) Healthcheck to use to validate VM health. Valid values: “INSTANCE_STATE”.
* `health_check_grace_period` - (Optional) Period of time (seconds) to wait for VM to reach healthiness before monitoring for unhealthiness.
* `auto_healing` - (Optional) Enable Autohealing of unhealthy VMs.

* `network` - (Required) Defines the Virtual Network and Subnet for your Elastigroup.
* `virtual_network_name` - (Required) Name of Vnet.
* `subnet_name` - (Required) ID of subnet.
* `resource_group_name` - (Required) Vnet Resource Group Name.
* `assign_public_up` - (Optional, Default: `false`) Assign a public IP to each VM in the Elastigroup.

* `user_name` - (Required) Set admin access for accessing your VMs.
* `ssh_public_key` - (Optional) SSH for admin access to Linux VMs. Required for Linux product types.
* `password` - (Optional) Password for admin access to Windows VMs. Required for Windows product types.

