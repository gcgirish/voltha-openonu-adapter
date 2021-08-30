/*
 * Copyright (c) 2018 - present.  Boling Consulting Solutions (bcsw.net)
 * Copyright 2020-present Open Networking Foundation

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 * http://www.apache.org/licenses/LICENSE-2.0

 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/*
 * NOTE: This file was generated, manual edits will be overwritten!
 *
 * Generated by 'goCodeGenerator.py':
 *              https://github.com/cboling/OMCI-parser/README.md
 */

package generated

import "github.com/deckarep/golang-set"

// Onu2GClassID is the 16-bit ID for the OMCI
// Managed entity ONU2-G
const Onu2GClassID ClassID = ClassID(257)

var onu2gBME *ManagedEntityDefinition

// Onu2G (class ID #257)
//	This ME contains additional attributes associated with a PON ONU. The ONU automatically creates
//	an instance of this ME. Its attributes are populated according to data within the ONU itself.
//
//	This ME is the same as the ONT2-G of [ITUT G.984.4], with extensions.
//
//	Relationships
//		This ME is paired with the ONU-G entity.
//
//	Attributes
//		Managed Entity Id
//			Managed entity ID: This attribute uniquely identifies each instance of this ME. There is only
//			one instance, number 0. (R) (mandatory) (2-bytes)
//
//		Equipment Id
//			Equipment ID: This attribute may be used to identify the specific type of ONU. In some
//			environments, this attribute may include the common language equipment identification (CLEI)
//			code. (R) (optional) (20-bytes)
//
//		Optical Network Unit Management And Control Channel Omcc Version
//			(R) (mandatory) (1-byte)
//
//		Vendor Product Code
//			Vendor product code: This attribute contains a vendor-specific product code for the ONU. (R)
//			(optional) (2-bytes)
//
//		Security Capability
//			(R) (mandatory) (1-byte)
//
//		Security Mode
//			Upon ME instantiation, the ONU sets this attribute to 1, AES-128. Attribute value 1 does not
//			imply that any channels are encrypted; that process is negotiated at the PLOAM layer. It only
//			signifies that the advanced encryption standard (AES) with 128-bit keys is the security mode to
//			be used on any channels that the OLT may choose to encrypt. (R,-W) (mandatory) (1-byte)
//
//		Total Priority Queue Number
//			Total priority queue number: This attribute reports the total number of upstream priority queues
//			that are not associated with a circuit pack, but with the ONU in its entirety. Upon ME
//			instantiation, the ONU sets this attribute to the value that represents its capabilities. (R)
//			(mandatory) (2-bytes)
//
//		Total Traffic Scheduler Number
//			Total traffic scheduler number: This attribute reports the total number of traffic schedulers
//			that are not associated with a circuit pack, but with the ONU in its entirety. The ONU supports
//			null function, strict priority scheduling and weighted round robin (WRR) from the priority
//			control and guarantee of minimum rate control points of view, respectively. If the ONU has no
//			global traffic schedulers, this attribute is 0. (R) (mandatory) (1-byte)
//
//		Deprecated
//			Deprecated:	This attribute should always be set to 1 by the ONU and ignored by the OLT. (R)
//			(mandatory) (1-byte)
//
//		Total Gem Port_Id Number
//			Total GEM port-ID number: This attribute reports the total number of GEM port-IDs supported by
//			the ONU. The maximum value is specified in the corresponding TC recommendations. Upon ME
//			instantiation, the ONU sets this attribute to the value that represents its capabilities. (R)
//			(optional) (2-bytes)
//
//		Sysuptime
//			SysUpTime:	This attribute counts 10 ms intervals since the ONU was last initialized. It rolls
//			over to 0 when full (see [IETF RFC 1213]). (R) (optional) (4-bytes)
//
//		Connectivity Capability
//			(R) (optional) (2 bytes)
//
//		Current Connectivity Mode
//			(R, W) (optional) (1 byte)
//
//		Quality Of Service Qos Configuration Flexibility
//			The ME ID of both the T-CONT and traffic scheduler contains a slot number. Even when attributes
//			in the above list are RW, it is never permitted to change the slot number in a reference. That
//			is, configuration flexibility never extends across slots. It is also not permitted to change the
//			directionality of an upstream queue to downstream or vice versa.
//
//		Priority Queue Scale Factor
//			NOTE 3 - Some legacy implementations may take the queue scale factor from the GEM block length
//			attribute of the ANI-G ME. That option is discouraged in new implementations.
//
type Onu2G struct {
	ManagedEntityDefinition
	Attributes AttributeValueMap
}

func init() {
	onu2gBME = &ManagedEntityDefinition{
		Name:    "Onu2G",
		ClassID: 257,
		MessageTypes: mapset.NewSetWith(
			Get,
			Set,
		),
		AllowedAttributeMask: 0xfffc,
		AttributeDefinitions: AttributeDefinitionMap{
			0:  Uint16Field("ManagedEntityId", PointerAttributeType, 0x0000, 0, mapset.NewSetWith(Read), false, false, false, 0),
			1:  MultiByteField("EquipmentId", StringAttributeType, 0x8000, 20, toOctets("AAAAAAAAAAAAAAAAAAAAAAAAAAA="), mapset.NewSetWith(Read), false, true, false, 1),
			2:  ByteField("OpticalNetworkUnitManagementAndControlChannelOmccVersion", EnumerationAttributeType, 0x4000, 164, mapset.NewSetWith(Read), true, false, false, 2),
			3:  Uint16Field("VendorProductCode", UnsignedIntegerAttributeType, 0x2000, 0, mapset.NewSetWith(Read), false, true, false, 3),
			4:  ByteField("SecurityCapability", EnumerationAttributeType, 0x1000, 0, mapset.NewSetWith(Read), false, false, false, 4),
			5:  ByteField("SecurityMode", EnumerationAttributeType, 0x0800, 0, mapset.NewSetWith(Read, Write), false, false, false, 5),
			6:  Uint16Field("TotalPriorityQueueNumber", UnsignedIntegerAttributeType, 0x0400, 0, mapset.NewSetWith(Read), false, false, false, 6),
			7:  ByteField("TotalTrafficSchedulerNumber", UnsignedIntegerAttributeType, 0x0200, 0, mapset.NewSetWith(Read), false, false, false, 7),
			8:  ByteField("Deprecated", UnsignedIntegerAttributeType, 0x0100, 1, mapset.NewSetWith(Read), false, false, true, 8),
			9:  Uint16Field("TotalGemPortIdNumber", UnsignedIntegerAttributeType, 0x0080, 0, mapset.NewSetWith(Read), false, true, false, 9),
			10: Uint32Field("Sysuptime", UnsignedIntegerAttributeType, 0x0040, 0, mapset.NewSetWith(Read), false, true, false, 10),
			11: Uint16Field("ConnectivityCapability", BitFieldAttributeType, 0x0020, 0, mapset.NewSetWith(Read), false, true, false, 11),
			12: ByteField("CurrentConnectivityMode", BitFieldAttributeType, 0x0010, 0, mapset.NewSetWith(Read, Write), false, true, false, 12),
			13: Uint16Field("QualityOfServiceQosConfigurationFlexibility", BitFieldAttributeType, 0x0008, 0, mapset.NewSetWith(Read), false, true, false, 13),
			14: Uint16Field("PriorityQueueScaleFactor", UnsignedIntegerAttributeType, 0x0004, 1, mapset.NewSetWith(Read, Write), false, true, false, 14),
		},
		Access:  CreatedByOnu,
		Support: UnknownSupport,
	}
}

// NewOnu2G (class ID 257) creates the basic
// Managed Entity definition that is used to validate an ME of this type that
// is received from or transmitted to the OMCC.
func NewOnu2G(params ...ParamData) (*ManagedEntity, OmciErrors) {
	return NewManagedEntity(*onu2gBME, params...)
}