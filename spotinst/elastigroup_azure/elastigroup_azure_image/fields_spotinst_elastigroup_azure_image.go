package elastigroup_azure_image

import (
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/elastigroup/providers/azure"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/commons"
)

func Setup(fieldsMap map[commons.FieldName]*commons.GenericField) {

	fieldsMap[Image] = commons.NewGenericField(
		commons.ElastigroupAzureImage,
		Image,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(Marketplace): &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						//ConflictsWith: []string{string(Custom)},
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								string(Publisher): &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},

								string(Offer): &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},

								string(Sku): &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},

					string(Custom): &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						//ConflictsWith: []string{string(Marketplace)},
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								string(ImageName): &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},

								string(ResourceGroupName): &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},

		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			//egWrapper := resourceObject.(*commons.ElastigroupAzureWrapper)
			//elastigroup := egWrapper.GetElastigroup()
			//var value []interface{} = nil
			//
			//if elastigroup.Compute != nil && elastigroup.Compute.LaunchSpecification != nil &&
			//	elastigroup.Compute.LaunchSpecification.Image != nil {
			//
			//	if elastigroup.Compute.LaunchSpecification.Image.MarketPlace != nil {
			//		value = flattenAzureGroupMarketplaceImage(elastigroup.Compute.LaunchSpecification.Image.MarketPlace)
			//	}
			//
			//	if elastigroup.Compute.LaunchSpecification.Image.Custom != nil {
			//		value = flattenAzureGroupCustomImage(elastigroup.Compute.LaunchSpecification.Image.Custom)
			//	}
			//
			//	if err := resourceData.Set(string(Image), value); err != nil {
			//		return fmt.Errorf(string(commons.FailureFieldReadPattern), string(Image), err)
			//	}
			//
			//}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupAzureWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.GetOk(string(Image)); ok {
				if image, err := expandAzureGroupImage(v); err != nil {
					return err
				} else {
					elastigroup.Compute.LaunchSpecification.SetImage(image)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupAzureWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.GetOk(string(Image)); ok {
				if image, err := expandAzureGroupImage(v); err != nil {
					return err
				} else {
					elastigroup.Compute.LaunchSpecification.SetImage(image)
				}
			}
			return nil
		},
		nil,
	)
}

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Utils
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//func flattenAzureGroupImage(image *azure.Image) []interface{} {
//	result := make(map[string]interface{})
//	result[string(Marketplace)] = ""
//	result[string(Custom)] = ""
//}

//func flattenAWSGroupEphemeralBlockDevices(devices []*aws.BlockDeviceMapping) []interface{} {
//	result := make([]interface{}, 0, len(devices))
//	for _, dev := range devices {
//		if dev.EBS == nil {
//			m := make(map[string]interface{})
//			m[string(DeviceName)] = spotinst.StringValue(dev.DeviceName)
//			m[string(VirtualName)] = spotinst.StringValue(dev.VirtualName)
//			result = append(result, m)
//		}
//	}
//	return result
//}
//
//func flattenAzureGroupMarketplaceImage(image *azure.MarketPlaceImage) []interface{} {
//
//}
//
//func flattenAzureGroupCustomImage(image *azure.CustomImage) []interface{} {
//
//}

func expandAzureGroupImage(data interface{}) (*azure.Image, error) {
	image := &azure.Image{}
	list := data.([]interface{})
	if list != nil && list[0] != nil {
		m := list[0].(map[string]interface{})
		if v, ok := m[string(Marketplace)]; ok {
			marketplace, err := expandAzureGroupMarketplaceImage(v)
			if err != nil {
				return nil, err
			}

			if marketplace != nil {
				image.SetMarketPlaceImage(marketplace)
			}
		} else {
			image.MarketPlace = nil
		}

		if v, ok := m[string(Custom)]; ok {

			custom, err := expandAzureGroupCustomImage(v)
			if err != nil {
				return nil, err
			}
			if custom != nil {
				image.SetCustom(custom) // todo rename this to setcustomimage in sdk
			}
		} else {
			image.Custom = nil
		}

	} else {
		return nil, errors.New("invalid image configuration")
	}
	return image, nil
}

func expandAzureGroupMarketplaceImage(data interface{}) (*azure.MarketPlaceImage, error) {
	market := &azure.MarketPlaceImage{}
	if list := data.([]interface{}); len(list) > 0 {
		if list != nil && list[0] != nil {
			m := list[0].(map[string]interface{})

			if v, ok := m[string(Publisher)].(string); ok && v != "" {
				market.SetPublisher(spotinst.String(v))
			}

			if v, ok := m[string(Offer)].(string); ok && v != "" {
				market.SetOffer(spotinst.String(v))
			}

			if v, ok := m[string(Sku)].(string); ok && v != "" {
				market.SetSKU(spotinst.String(v))
			}

		}
		return market, nil
	}
	return nil, nil
}

func expandAzureGroupCustomImage(data interface{}) (*azure.CustomImage, error) {
	if list := data.([]interface{}); len(list) > 0 {
		custom := &azure.CustomImage{}
		if list != nil && list[0] != nil {
			m := list[0].(map[string]interface{})

			if v, ok := m[string(ImageName)].(string); ok && v != "" {
				custom.SetImageName(spotinst.String(v))
			}

			if v, ok := m[string(ResourceGroupName)].(string); ok && v != "" {
				custom.SetResourceGroupName(spotinst.String(v))
			}

		}
		return custom, nil
	}
	return nil, nil
}
