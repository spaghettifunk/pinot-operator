package sdk

import (
	"fmt"
	"os"

	models "github.com/spaghettifunk/pinot-go-client/models"
	operatorsv1alpha1 "github.com/spaghettifunk/pinot-operator/pkg/apis/pinot/v1alpha1"
	"github.com/spaghettifunk/pinot-operator/pkg/util"
)

const (
	pinotControllerHeadless = "pinot-controller-headless"
	kubernetesDomain        = "cluster.local"
)

// GeneratePinotControlleAddressWithoutPort returns the host of the controller service
func GeneratePinotControlleAddressWithoutPort(config *operatorsv1alpha1.Pinot) string {
	if os.Getenv("LOCAL_DEBUG") == "true" {
		return "localhost"
	}
	return fmt.Sprintf("%s.%s.svc.%s",
		pinotControllerHeadless,
		config.Namespace,
		kubernetesDomain,
	)
}

// GeneratePinotControllerAddress returns the fully qualified DNS of the controller server
func GeneratePinotControllerAddress(config *operatorsv1alpha1.Pinot) string {
	return fmt.Sprintf("%s:%d", GeneratePinotControlleAddressWithoutPort(config), config.Spec.Controller.Service.Port)
}

// ConvertCRDSchemaToSDKSchema will convert the crd object to an object that will be used by
// the SDK to submit the request
func ConvertCRDSchemaToSDKSchema(config *operatorsv1alpha1.Schema) *models.Schema {
	schema := &models.Schema{
		SchemaName:          config.Spec.Name,
		PrimaryKeyColumns:   config.Spec.PrimaryKeys,
		DimensionFieldSpecs: convertDimensionFieldSpec(config.Spec.Dimensions),
		MetricFieldSpecs:    convertMetricFieldSpec(config.Spec.Metrics),
		DateTimeFieldSpecs:  convertDateTimeFieldSpec(config.Spec.DateTimes),
		TimeFieldSpec:       convertTimeFieldSpec(config.Spec.TimeField),
	}
	return schema
}

func convertDimensionFieldSpec(dimensions []*operatorsv1alpha1.DimensionFieldSpec) []*models.DimensionFieldSpec {
	list := make([]*models.DimensionFieldSpec, len(dimensions))
	for _, dim := range dimensions {
		list = append(list, &models.DimensionFieldSpec{
			Name:                   util.PointerToString(dim.Name),
			DataType:               util.PointerToString(dim.DataType),
			DefaultNullValue:       dim.DefaultNullValue,
			DefaultNullValueString: util.PointerToString(dim.DefaultNullValueString),
			MaxLength:              util.PointerToInt32(dim.MaxLength),
			SingleValueField:       util.PointerToBool(dim.SingleValueField),
			TransformFunction:      util.PointerToString(dim.TransformFunction),
			VirtualColumnProvider:  util.PointerToString(dim.VirtualColumnProvider),
		})
	}
	return list
}

func convertMetricFieldSpec(metrics []*operatorsv1alpha1.MetricFieldSpec) []*models.MetricFieldSpec {
	list := make([]*models.MetricFieldSpec, len(metrics))
	for _, mt := range metrics {
		list = append(list, &models.MetricFieldSpec{
			Name:                   util.PointerToString(mt.Name),
			DataType:               util.PointerToString(mt.DataType),
			DefaultNullValue:       mt.DefaultNullValue,
			DefaultNullValueString: util.PointerToString(mt.DefaultNullValueString),
			MaxLength:              util.PointerToInt32(mt.MaxLength),
			SingleValueField:       util.PointerToBool(mt.SingleValueField),
			TransformFunction:      util.PointerToString(mt.TransformFunction),
			VirtualColumnProvider:  util.PointerToString(mt.VirtualColumnProvider),
		})
	}
	return list
}

func convertDateTimeFieldSpec(datetimes []*operatorsv1alpha1.DatetimeFieldSpec) []*models.DateTimeFieldSpec {
	list := make([]*models.DateTimeFieldSpec, len(datetimes))
	for _, dt := range datetimes {
		list = append(list, &models.DateTimeFieldSpec{
			Name:                   util.PointerToString(dt.Name),
			DataType:               util.PointerToString(dt.DataType),
			DefaultNullValue:       dt.DefaultNullValue,
			DefaultNullValueString: util.PointerToString(dt.DefaultNullValueString),
			MaxLength:              util.PointerToInt32(dt.MaxLength),
			SingleValueField:       util.PointerToBool(dt.SingleValueField),
			TransformFunction:      util.PointerToString(dt.TransformFunction),
			VirtualColumnProvider:  util.PointerToString(dt.VirtualColumnProvider),
			Format:                 util.PointerToString(dt.Format),
			Granularity:            util.PointerToString(dt.Granularity),
		})
	}
	return list
}

func convertTimeFieldSpec(tf *operatorsv1alpha1.TimeFieldSpec) *models.TimeFieldSpec {
	timefield := &models.TimeFieldSpec{
		Name:                    util.PointerToString(tf.Name),
		DataType:                util.PointerToString(tf.DataType),
		DefaultNullValue:        tf.DefaultNullValue,
		DefaultNullValueString:  util.PointerToString(tf.DefaultNullValueString),
		MaxLength:               util.PointerToInt32(tf.MaxLength),
		SingleValueField:        util.PointerToBool(tf.SingleValueField),
		TransformFunction:       util.PointerToString(tf.TransformFunction),
		VirtualColumnProvider:   util.PointerToString(tf.VirtualColumnProvider),
		IncomingGranularitySpec: &models.TimeGranularitySpec{},
		OutgoingGranularitySpec: &models.TimeGranularitySpec{},
	}
	if tf.IncomingGranularity != nil {
		timefield.IncomingGranularitySpec.DataType = util.PointerToString(tf.IncomingGranularity.DataType)
		timefield.IncomingGranularitySpec.Name = util.PointerToString(tf.IncomingGranularity.Name)
		timefield.IncomingGranularitySpec.TimeFormat = util.PointerToString(tf.IncomingGranularity.TimeFormat)
		timefield.IncomingGranularitySpec.TimeType = util.PointerToString(tf.IncomingGranularity.TimeType)
		timefield.IncomingGranularitySpec.TimeUnitSize = util.PointerToInt32(tf.IncomingGranularity.TimeUnitSize)
	}
	if tf.OutgoingGranularity != nil {
		timefield.OutgoingGranularitySpec.DataType = util.PointerToString(tf.OutgoingGranularity.DataType)
		timefield.OutgoingGranularitySpec.Name = util.PointerToString(tf.OutgoingGranularity.Name)
		timefield.OutgoingGranularitySpec.TimeFormat = util.PointerToString(tf.OutgoingGranularity.TimeFormat)
		timefield.OutgoingGranularitySpec.TimeType = util.PointerToString(tf.OutgoingGranularity.TimeType)
		timefield.OutgoingGranularitySpec.TimeUnitSize = util.PointerToInt32(tf.OutgoingGranularity.TimeUnitSize)
	}
	return timefield
}
