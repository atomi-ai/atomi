package models

import (
	"time"
)

type TokenResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type QuoteRequest struct {
	// Street Address, City, State, Zip
	DropoffAddress     string     `json:"dropoff_address"`
	PickupAddress      string     `json:"pickup_address"`
	DropoffLatitude    *float64   `json:"dropoff_latitude,omitempty"`
	DropoffLongitude   *float64   `json:"dropoff_longitude,omitempty"`
	DropoffPhoneNumber *string    `json:"dropoff_phone_number,omitempty"`
	PickupLatitude     *float64   `json:"pickup_latitude,omitempty"`
	PickupLongitude    *float64   `json:"pickup_longitude,omitempty"`
	PickupPhoneNumber  *string    `json:"pickup_phone_number,omitempty"`
	PickupReadyDt      *time.Time `json:"pickup_ready_dt,omitempty"`
	PickupDeadlineDt   *time.Time `json:"pickup_deadline_dt,omitempty"`
	DropoffReadyDt     *time.Time `json:"dropoff_ready_dt,omitempty"`
	DropoffDeadlineDt  *time.Time `json:"dropoff_deadline_dt,omitempty"`
	ManifestTotalValue *int64     `json:"manifest_total_value,omitempty"`
	ExternalStoreID    *string    `json:"external_store_id,omitempty"`
}

type QuoteResponse struct {
	Created         time.Time `json:"created"`
	CurrencyType    string    `json:"currency_type"`
	DropoffDeadline time.Time `json:"dropoff_deadline"`
	DropoffETA      time.Time `json:"dropoff_eta"`
	Duration        int64     `json:"duration"`
	Expires         time.Time `json:"expires"`
	Fee             int64     `json:"fee"`
	ID              string    `json:"id"`
	Kind            string    `json:"kind"`
	PickupDuration  int64     `json:"pickup_duration"`
	ExternalStoreID *string   `json:"external_store_id,omitempty"`
}

type DeliveryData struct {
	DropoffAddress      string                   `json:"dropoff_address"`
	DropoffName         string                   `json:"dropoff_name"`
	DropoffPhoneNumber  string                   `json:"dropoff_phone_number"`
	ManifestItems       []ManifestItem           `json:"manifest_items"`
	PickupAddress       string                   `json:"pickup_address"`
	PickupName          string                   `json:"pickup_name"`
	PickupPhoneNumber   string                   `json:"pickup_phone_number"`
	DeliverableAction   *DeliverableAction       `json:"deliverable_action,omitempty"`
	DropoffBusinessName *string                  `json:"dropoff_business_name,omitempty"`
	DropoffLatitude     *float64                 `json:"dropoff_latitude,omitempty"`
	DropoffLongitude    *float64                 `json:"dropoff_longitude,omitempty"`
	DropoffNotes        *string                  `json:"dropoff_notes,omitempty"`
	DropoffSellerNotes  *string                  `json:"dropoff_seller_notes,omitempty"`
	DropoffVerification *VerificationRequirement `json:"dropoff_verification,omitempty"`
	ManifestReference   *string                  `json:"manifest_reference,omitempty"`
	ManifestTotalValue  *int                     `json:"manifest_total_value,omitempty"`
	PickupBusinessName  *string                  `json:"pickup_business_name,omitempty"`
	PickupLatitude      *float64                 `json:"pickup_latitude,omitempty"`
	PickupLongitude     *float64                 `json:"pickup_longitude,omitempty"`
	PickupNotes         *string                  `json:"pickup_notes,omitempty"`
	PickupVerification  *VerificationRequirement `json:"pickup_verification,omitempty"`
	QuoteID             *string                  `json:"quote_id,omitempty"`
	UndeliverableAction *UndeliverableAction     `json:"undeliverable_action,omitempty"`
	PickupReadyDt       *time.Time               `json:"pickup_ready_dt,omitempty"`
	PickupDeadlineDt    *time.Time               `json:"pickup_deadline_dt,omitempty"`
	DropoffReadyDt      *time.Time               `json:"dropoff_ready_dt,omitempty"`
	DropoffDeadlineDt   *time.Time               `json:"dropoff_deadline_dt,omitempty"`
	RequiresID          *bool                    `json:"requires_id,omitempty"`
	Tip                 *int                     `json:"tip,omitempty"`
	IdempotencyKey      *string                  `json:"idempotency_key,omitempty"`
	ExternalStoreID     *string                  `json:"external_store_id,omitempty"`
	ReturnVerification  *VerificationRequirement `json:"return_verification,omitempty"`
	TestSpecifications  *TestSpecifications      `json:"test_specifications,omitempty"`
}

type Size string

const (
	SizeSmall  Size = "small"
	SizeMedium Size = "medium"
	SizeLarge  Size = "large"
	SizeXLarge Size = "xlarge"
)

type Dimensions struct {
	Length *float64 `json:"length,omitempty"`
	Height *float64 `json:"height,omitempty"`
	Depth  *float64 `json:"depth,omitempty"`
}

type ManifestItem struct {
	Name       string      `json:"name"`
	Quantity   int         `json:"quantity"`
	Size       *Size       `json:"size"`
	Dimensions *Dimensions `json:"dimensions,omitempty"`
	Price      *int        `json:"price,omitempty"`
	Weight     *int        `json:"weight,omitempty"`
}

type DeliverableAction string

const (
	DeliverableActionMeetAtDoor  DeliverableAction = "deliverable_action_meet_at_door"
	DeliverableActionLeaveAtDoor DeliverableAction = "deliverable_action_leave_at_door"
)

type SignatureRequirement struct {
	Enabled                   bool  `json:"enabled"`
	CollectSignerName         *bool `json:"collect_signer_name,omitempty"`
	CollectSignerRelationship *bool `json:"collect_signer_relationship,omitempty"`
}

type BarcodeRequirement struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type PincodeRequirement struct {
	Enabled bool    `json:"enabled"`
	Value   *string `json:"value,omitempty"`
}

type PackageRequirement struct {
	BagCount   *int `json:"bag_count,omitempty"`
	DrinkCount *int `json:"drink_count,omitempty"`
}

type IdentificationRequirement struct {
	MinAge *int `json:"min_age,omitempty"`
}

type VerificationRequirement struct {
	SignatureRequirement *SignatureRequirement      `json:"signature_requirement,omitempty"`
	Barcodes             []BarcodeRequirement       `json:"barcodes,omitempty"`
	Pincode              *PincodeRequirement        `json:"pincode,omitempty"`
	Package              *PackageRequirement        `json:"package,omitempty"`
	Identification       *IdentificationRequirement `json:"identification,omitempty"`
	Picture              *bool                      `json:"picture,omitempty"`
}

type UndeliverableAction string

const (
	UndeliverableActionLeaveAtDoor UndeliverableAction = "leave_at_door"
	UndeliverableActionReturn      UndeliverableAction = "return"
)

type RoboCourierSpecification struct {
	Mode               string                   `json:"mode"`
	EnrouteForPickupAt *string                  `json:"enroute_for_pickup_at,omitempty"`
	PickupImminentAt   *string                  `json:"pickup_imminent_at,omitempty"`
	PickupAt           *string                  `json:"pickup_at,omitempty"`
	DropoffImminentAt  *string                  `json:"dropoff_imminent_at,omitempty"`
	DropoffAt          *string                  `json:"dropoff_at,omitempty"`
	CancelReason       *RoboCourierCancelReason `json:"cancel_reason,omitempty"`
}

type RoboCourierCancelReason string

const (
	CannotAccessCustomerLocation RoboCourierCancelReason = "cannot_access_customer_location"
	CannotFindCustomerAddress    RoboCourierCancelReason = "cannot_find_customer_address"
	CustomerRejectedOrder        RoboCourierCancelReason = "customer_rejected_order"
	CustomerUnavailable          RoboCourierCancelReason = "customer_unavailable"
)

type TestSpecifications struct {
	RoboCourierSpecification RoboCourierSpecification `json:"robo_courier_specification"`
}

type DeliveryResponse struct {
	Complete            bool               `json:"complete"`
	Courier             *CourierInfo       `json:"courier,omitempty"`
	CourierImminent     bool               `json:"courier_imminent"`
	Created             string             `json:"created"`
	Currency            string             `json:"currency"`
	Dropoff             WaypointInfo       `json:"dropoff"`
	DropoffDeadline     string             `json:"dropoff_deadline"`
	DropoffETA          string             `json:"dropoff_eta"`
	DropoffIdentifier   *string            `json:"dropoff_identifier,omitempty"`
	DropoffReady        string             `json:"dropoff_ready"`
	ExternalID          string             `json:"external_id"`
	Fee                 int                `json:"fee"`
	ID                  string             `json:"id"`
	Kind                string             `json:"kind"`
	LiveMode            bool               `json:"live_mode"`
	Manifest            ManifestInfo       `json:"manifest"`
	ManifestItems       []ManifestItem     `json:"manifest_items"`
	Pickup              WaypointInfo       `json:"pickup"`
	PickupDeadline      string             `json:"pickup_deadline"`
	PickupETA           string             `json:"pickup_eta"`
	PickupReady         string             `json:"pickup_ready"`
	QuoteID             string             `json:"quote_id"`
	Refund              *[]RefundData      `json:"refund,omitempty"`
	RelatedDeliveries   *[]RelatedDelivery `json:"related_deliveries,omitempty"`
	Status              DeliveryStatus     `json:"status"`
	Tip                 *int               `json:"tip,omitempty"`
	TrackingURL         string             `json:"tracking_url"`
	UndeliverableAction string             `json:"undeliverable_action"`
	UndeliverableReason string             `json:"undeliverable_reason"`
	Updated             string             `json:"updated"`
	UUID                string             `json:"uuid"`
	Return              *WaypointInfo      `json:"return,omitempty"`
}

type CourierInfo struct {
	Name        string   `json:"name"`
	Rating      *float64 `json:"rating,omitempty"`
	VehicleType string   `json:"vehicle_type"`
	PhoneNumber string   `json:"phone_number"`
	Location    LatLng   `json:"location"`
	ImgHref     *string  `json:"img_href,omitempty"`
}

type WaypointInfo struct {
	Name                     string                   `json:"name"`
	PhoneNumber              string                   `json:"phone_number"`
	Address                  string                   `json:"address"`
	DetailedAddress          *Address                 `json:"detailed_address,omitempty"`
	Notes                    *string                  `json:"notes,omitempty"`
	SellerNotes              *string                  `json:"seller_notes,omitempty"`
	CourierNotes             *string                  `json:"courier_notes,omitempty"`
	Location                 LatLng                   `json:"location"`
	Verification             *VerificationProof       `json:"verification,omitempty"`
	VerificationRequirements *VerificationRequirement `json:"verification_requirements,omitempty"`
	ExternalStoreID          *string                  `json:"external_store_id,omitempty"`
}

type ManifestInfo struct {
	Reference   string  `json:"reference"`
	Description *string `json:"description,omitempty"`
	TotalValue  int     `json:"total_value"`
}

type RefundData struct {
	ID                 string            `json:"id"`
	CreatedAt          int               `json:"created_at"`
	CurrencyCode       string            `json:"currency_code"`
	TotalPartnerRefund int               `json:"total_partner_refund"`
	TotalUberRefund    int               `json:"total_uber_refund"`
	RefundFees         []RefundFee       `json:"refund_fees"`
	RefundOrderItems   []RefundOrderItem `json:"refund_order_items"`
}

type RelatedDelivery struct {
	ID           string `json:"id"`
	Relationship string `json:"relationship"`
}

type DeliveryStatus string

const (
	DeliveryStatusPending        DeliveryStatus = "pending"
	DeliveryStatusPickup         DeliveryStatus = "pickup"
	DeliveryStatusPickupComplete DeliveryStatus = "pickup_complete"
	DeliveryStatusDropoff        DeliveryStatus = "dropoff"
	DeliveryStatusDelivered      DeliveryStatus = "delivered"
	DeliveryStatusCanceled       DeliveryStatus = "canceled"
	DeliveryStatusReturned       DeliveryStatus = "returned"
)

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type VerificationProof struct {
	Signature          *SignatureProof      `json:"signature,omitempty"`
	Barcodes           []BarcodeRequirement `json:"barcodes,omitempty"`
	Picture            *PictureProof        `json:"picture,omitempty"`
	Identification     *IdentificationProof `json:"identification,omitempty"`
	PinCode            *PincodeProof        `json:"pin_code,omitempty"`
	CompletionLocation *LatLng              `json:"completion_location,omitempty"`
}

type RefundFee struct {
	FeeCode  FeeCode  `json:"fee_code"`
	Value    int      `json:"value"`
	Category Category `json:"category"`
}

type RefundOrderItem struct {
	RefundItems         []RefundItem `json:"refund_items"`
	PartyAtFault        string       `json:"party_at_fault"`
	PartnerRefundAmount int          `json:"partner_refund_amount"`
	UberRefundAmount    int          `json:"uber_refund_amount"`
	Reason              string       `json:"reason"`
}

type SignatureProof struct {
	ImageURL           string  `json:"image_url"`
	SignerName         *string `json:"signer_name,omitempty"`
	SignerRelationship *string `json:"signer_relationship,omitempty"`
}

type PictureProof struct {
	ImageURL string `json:"image_url"`
}

type IdentificationProof struct {
	MinAgeVerified *bool `json:"min_age_verified,omitempty"`
}

type PincodeProof struct {
	Entered string `json:"entered"`
}

type FeeCode string

const (
	FeeCodeUberDeliveryFee = FeeCode("UBER_DELIVERY_FEE")
	FeeCodePartnerFee      = FeeCode("PARTNER_FEE")
	FeeCodePartnerTax      = FeeCode("PARTNER_TAX")
)

type Category string

const (
	CategoryDelivery = Category("DELIVERY")
	CategoryTax      = Category("TAX")
)

type RefundItem struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}
