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

// EfmBondingLinkClassID is the 16-bit ID for the OMCI
// Managed entity EFM bonding link
const EfmBondingLinkClassID ClassID = ClassID(419)

var efmbondinglinkBME *ManagedEntityDefinition

// EfmBondingLink (class ID #419)
//	The EFM bonding link represents a link that can be bonded with other links to form a group. In
//	[IEEE 802.3], a bonding group is known as a PAF and a link is known as a PME. Instances of this
//	ME are created and deleted by the OLT.
//
//	Relationships
//		An instance of this ME may be associated with zero or one instance of an EFM bonding group.
//
//	Attributes
//		Managed Entity Id
//			NOTE - This attribute has the same meaning as the Stream ID in clause C.3.1.2 of [ITU-T
//			G.998.2], except that it cannot be changed. (R, setbycreate) (mandatory) (2-bytes)
//
//		Associated Group Me Id
//			Associated group ME ID: This attribute is the ME ID of the bonding group to which this link is
//			associated. Changing this attribute moves the link from one group to another. Setting this
//			attribute to an ME ID that has not yet been provisioned will result in this link being placed in
//			a single-link group that contains only this link. The default value for this attribute is the
//			null pointer, 0xFFFF. (R,-W, setbycreate) (mandatory) (2-bytes)
//
//		Link Alarm Enable
//			(R,-W, setbycreate) (mandatory) (1-bytes)
//
type EfmBondingLink struct {
	ManagedEntityDefinition
	Attributes AttributeValueMap
}

func init() {
	efmbondinglinkBME = &ManagedEntityDefinition{
		Name:    "EfmBondingLink",
		ClassID: 419,
		MessageTypes: mapset.NewSetWith(
			Create,
			Delete,
			Get,
			Set,
		),
		AllowedAttributeMask: 0xc000,
		AttributeDefinitions: AttributeDefinitionMap{
			0: Uint16Field("ManagedEntityId", PointerAttributeType, 0x0000, 0, mapset.NewSetWith(Read, SetByCreate), false, false, false, 0),
			1: Uint16Field("AssociatedGroupMeId", UnsignedIntegerAttributeType, 0x8000, 0, mapset.NewSetWith(Read, SetByCreate, Write), false, false, false, 1),
			2: ByteField("LinkAlarmEnable", UnsignedIntegerAttributeType, 0x4000, 0, mapset.NewSetWith(Read, SetByCreate, Write), false, false, false, 2),
		},
		Access:  CreatedByOlt,
		Support: UnknownSupport,
		Alarms: AlarmMap{
			0: "Link down",
		},
	}
}

// NewEfmBondingLink (class ID 419) creates the basic
// Managed Entity definition that is used to validate an ME of this type that
// is received from or transmitted to the OMCC.
func NewEfmBondingLink(params ...ParamData) (*ManagedEntity, OmciErrors) {
	return NewManagedEntity(*efmbondinglinkBME, params...)
}
