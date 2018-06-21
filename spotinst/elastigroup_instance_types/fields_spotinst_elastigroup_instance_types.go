package elastigroup_instance_types

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/elastigroup/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/commons"
)

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Setup
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
func Setup(fieldsMap map[commons.FieldName]*commons.GenericField) {

	fieldsMap[OnDemand] = commons.NewGenericField(
		commons.ElastigroupInstanceType,
		OnDemand,
		&schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			elastigroup := resourceObject.(*aws.Group)
			var value *string = nil
			if elastigroup.Compute != nil && elastigroup.Compute.InstanceTypes != nil &&
				elastigroup.Compute.InstanceTypes.OnDemand != nil {
				value = elastigroup.Compute.InstanceTypes.OnDemand
			}
			if err := resourceData.Set(string(OnDemand), spotinst.StringValue(value)); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(OnDemand), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			elastigroup := resourceObject.(*aws.Group)
			if v, ok := resourceData.Get(string(OnDemand)).(string); ok && v != "" {
				elastigroup.Compute.InstanceTypes.SetOnDemand(spotinst.String(v))
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			elastigroup := resourceObject.(*aws.Group)
			if v, ok := resourceData.Get(string(OnDemand)).(string); ok && v != "" {
				elastigroup.Compute.InstanceTypes.SetOnDemand(spotinst.String(v))
			}
			return nil
		},
		nil,
	)

	fieldsMap[Spot] = commons.NewGenericField(
		commons.ElastigroupInstanceType,
		Spot,
		&schema.Schema{
			Type:     schema.TypeList,
			Required: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			elastigroup := resourceObject.(*aws.Group)
			if v, ok := resourceData.GetOk(string(Spot)); ok {
				spots := v.([]interface{})
				spotTypes := make([]string, len(spots))
				for i, j := range spots {
					spotTypes[i] = j.(string)
				}
				elastigroup.Compute.InstanceTypes.SetSpot(spotTypes)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			elastigroup := resourceObject.(*aws.Group)
			if v, ok := resourceData.GetOk(string(Spot)); ok {
				rawSpotTypes := v.([]interface{})
				spotTypes := make([]string, len(rawSpotTypes))
				for i, v := range rawSpotTypes {
					spotTypes[i] = v.(string)
				}
				elastigroup.Compute.InstanceTypes.SetSpot(spotTypes)
			}
			return nil
		},
		nil,
	)

	fieldsMap[InstanceTypeWeights] = commons.NewGenericField(
		commons.ElastigroupInstanceType,
		InstanceTypeWeights,
		&schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(InstanceType): &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},

					string(Weight): &schema.Schema{
						Type:     schema.TypeInt,
						Required: true,
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			elastigroup := resourceObject.(*aws.Group)
			if v, ok := resourceData.GetOkExists(string(InstanceTypeWeights)); ok && v != "" {
				if weights, err := expandAWSGroupInstanceTypeWeights(v); err != nil {
					return err
				} else {
					elastigroup.Compute.InstanceTypes.SetWeights(weights)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			elastigroup := resourceObject.(*aws.Group)
			var weightsToAdd []*aws.InstanceTypeWeight = nil
			if v, ok := resourceData.GetOk(string(InstanceTypeWeights)); ok {
				if weights, err := expandAWSGroupInstanceTypeWeights(v); err != nil {
					return err
				} else {
					weightsToAdd = weights
				}
			}
			elastigroup.Compute.InstanceTypes.SetWeights(weightsToAdd)
			return nil
		},
		nil,
	)
}

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Utils
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
func expandAWSGroupInstanceTypeWeights(data interface{}) ([]*aws.InstanceTypeWeight, error) {
	list := data.(*schema.Set).List()
	weights := make([]*aws.InstanceTypeWeight, 0, len(list))
	for _, v := range list {
		attr, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		if _, ok := attr[string(InstanceType)]; !ok {
			return nil, errors.New("[ERROR] Invalid instance type weight: instance_type missing")
		}

		if _, ok := attr[string(Weight)]; !ok {
			return nil, errors.New("[ERROR] Invalid instance type weight: weight missing")
		}
		instanceWeight := &aws.InstanceTypeWeight{}
		instanceWeight.SetInstanceType(spotinst.String(attr[string(InstanceType)].(string)))
		instanceWeight.SetWeight(spotinst.Int(attr[string(Weight)].(int)))
		//log.Printf("Group instance type weight configuration: %s", stringutil.Stringify(instanceWeight))
		weights = append(weights, instanceWeight)
	}
	return weights, nil
}