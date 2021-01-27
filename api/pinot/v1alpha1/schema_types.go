/*
Copyright 2021 the Apache Pinot Kubernetes Operator authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SchemaSpec defines the desired state of Schema
type SchemaSpec struct {
	// Name of the schema
	Name string `json:"name" protobuf:"byte,1,opt,name=name"`
	// PrimaryKeys is a list of columns that are set as primary keys
	PrimaryKeys []string `json:"primaryKeys" protobuf:"byte,2,opt,name=primaryKeys"`
	// Dimensions is a list of fields that represents the dimensions in the schema
	// ref: https://docs.pinot.apache.org/basics/components/schema#categories
	Dimensions []*DimensionFieldSpec `json:"dimensions" protobuf:"byte,3,opt,name=dimensions"`
	// Metrics is a list of fields that represents the metrics in the schema
	// ref: https://docs.pinot.apache.org/basics/components/schema#categories
	// +optional
	Metrics []*MetricFieldSpec `json:"metrics,omitempty" protobuf:"byte,4,opt,name=metrics"`
	// DateTimes is a list of fields that represents the datetimes in the schema
	// ref: https://docs.pinot.apache.org/basics/components/schema#categories
	DateTimes []*DatetimeFieldSpec `json:"dateTimes" protobuf:"byte,5,opt,name=dateTimes"`
	// TimeField represents the granularity
	TimeField *TimeFieldSpec `json:"timeField" protobuf:"byte,6,opt,name=timeField"`
	// +optional
	PinotServer *NamespacedName `json:"pinotServer,omitempty"`
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
}

// DimensionFieldSpec is typically used in slice and dice operations for answering business queries
type DimensionFieldSpec struct {
	CommonColumnSpec
}

// MetricFieldSpec represents the quantitative data of the table. Such columns are used for aggregation.
// In data warehouse terminology, these can also be referred to as fact or measure columns
type MetricFieldSpec struct {
	CommonColumnSpec
}

// DatetimeFieldSpec represents time columns in the data. There can be multiple time columns in a table, but only one of them
// can be treated as primary. Primary time column is the one that is present in the segment config.
type DatetimeFieldSpec struct {
	Format      *string `json:"format"`
	Granularity *string `json:"granularity"`
	CommonColumnSpec
}

// TimeFieldSpec represents the granularity for both ingestion and query segments
type TimeFieldSpec struct {
	// +optional
	IncomingGranularity *TimeGranularitySpec `json:"incomingGranularity"`
	// +optional
	OutgoingGranularity *TimeGranularitySpec `json:"outgoingGranularity"`
	// +optional
	CommonColumnSpec
}

// TimeGranularitySpec represents the granularity object
type TimeGranularitySpec struct {
	// Name of the time granularity specification
	Name *string `json:"name"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=INT;LONG;FLOAT;DOUBLE;BOOLEAN;STRING;BYTES;STRUCT;MAP;LIST
	// +kubebuilder:default:=none
	DataType *string `json:"dataType"`
	// TimeType is one of  TimeUnit enum values. e.g. HOURS , MINUTES etc. If your date is not in EPOCH format,
	// this value is not used and can be set to MILLISECONDS or any other unit.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=NANOSECONDS;MICROSECONDS;MILLISECONDS;SECONDS;MINUTES;HOURS;DAYS
	// +kubebuilder:default:=none
	TimeType *string `json:"typeType"`
	// TimeUnitSize is multiplied to the value present in the time column to get an actual timestamp.
	// eg: if timesize is 5 and value in time column is 4996308 minutes. The value that will be converted
	// to epoch timestamp will be 4996308 * 5 * 60 * 1000 = 1498892400000 milliseconds.
	// If your date is not in EPOCH format, this value is not used and can be set to 1 or any other integer.
	TimeUnitSize *int32 `json:"timeUnitSize"`
	// TimeFormat can be either EPOCH or SIMPLE_DATE_FORMAT. If it is SIMPLE_DATE_FORMAT, the pattern string is also specified.
	//
	// Here are some sample date-time formats you can use in the schema:
	// 1:MILLISECONDS:EPOCH - used when timestamp is in the epoch milliseconds and stored in LONG format
	// 1:HOURS:EPOCH - used when timestamp is in the epoch hours and stored in LONG  or INT format
	// 1:DAYS:SIMPLE_DATE_FORMAT:yyyy-MM-dd - when date is in STRING format and has the pattern year-month-date. e.g. 2020-08-21
	// 1:HOURS:SIMPLE_DATE_FORMAT:EEE MMM dd HH:mm:ss ZZZ yyyy - when date is in STRING format. e.g. s Mon Aug 24 12:36:50 America/Los_Angeles 2019
	TimeFormat *string `json:"timeFormat"`
}

// CommonColumnSpec represents the default values of such column
type CommonColumnSpec struct {
	Name             *string `json:"name"`
	SingleValueField *bool   `json:"singleValueField"`
	MaxLength        *int32  `json:"maxLength"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=INT;LONG;FLOAT;DOUBLE;BOOLEAN;STRING;BYTES;STRUCT;MAP;LIST
	// +kubebuilder:default:=none
	DataType               *string     `json:"dataType"`
	DefaultNullValue       interface{} `json:"defaultNullValue"`
	VirtualColumnProvider  *string     `json:"virtualColumnProvider"`
	TransformFunction      *string     `json:"transformFunction"`
	DefaultNullValueString *string     `json:"defaultNullValueString"`
}

// SchemaStatus defines the observed state of Schema
type SchemaStatus struct {
	Status       ConfigState `json:"Status,omitempty"`
	ErrorMessage string      `json:"ErrorMessage,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Schema is the Schema for the schemas API
// +k8s:openapi-gen=true
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:printcolumn:name="Error",type="string",JSONPath=".status.ErrorMessage",description="Error message"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Pinot Cluster",type="string",JSONPath=".spec.pinotServer"
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=schemas,shortName=scs
type Schema struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SchemaSpec   `json:"spec,omitempty"`
	Status SchemaStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SchemaList contains a list of Schema
type SchemaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Schema `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Schema{}, &SchemaList{})
}
