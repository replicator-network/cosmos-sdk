package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	commitment "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment"
	host "github.com/cosmos/cosmos-sdk/x/ibc/24-host"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
)

var _ sdk.Msg = MsgChannelOpenInit{}

type MsgChannelOpenInit struct {
	PortID    string         `json:"port_id"`
	ChannelID string         `json:"channel_id"`
	Channel   Channel        `json:"channel"`
	Signer    sdk.AccAddress `json:"signer"`
}

// NewMsgChannelOpenInit creates a new MsgChannelCloseInit MsgChannelOpenInit
func NewMsgChannelOpenInit(
	portID, channelID string, version string, channelOrder Order, connectionHops []string,
	counterpartyPortID, counterpartyChannelID string, signer sdk.AccAddress,
) MsgChannelOpenInit {
	counterparty := NewCounterparty(counterpartyPortID, counterpartyChannelID)
	channel := NewChannel(INIT, channelOrder, counterparty, connectionHops, version)
	return MsgChannelOpenInit{
		PortID:    portID,
		ChannelID: channelID,
		Channel:   channel,
		Signer:    signer,
	}
}

// Route implements sdk.Msg
func (msg MsgChannelOpenInit) Route() string {
	return ibctypes.RouterKey
}

// Type implements sdk.Msg
func (msg MsgChannelOpenInit) Type() string {
	return "channel_open_init"
}

// ValidateBasic implements sdk.Msg
func (msg MsgChannelOpenInit) ValidateBasic() sdk.Error {
	if err := host.DefaultPortIdentifierValidator(msg.PortID); err != nil {
		return sdk.ConvertError(sdkerrors.Wrap(err, "invalid port ID"))
	}
	if err := host.DefaultClientIdentifierValidator(msg.ChannelID); err != nil {
		return sdk.ConvertError(sdkerrors.Wrap(err, "invalid channel ID"))
	}
	// Signer can be empty
	return sdk.ConvertError(msg.Channel.ValidateBasic())
}

// GetSignBytes implements sdk.Msg
func (msg MsgChannelOpenInit) GetSignBytes() []byte {
	return sdk.MustSortJSON(SubModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgChannelOpenInit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

var _ sdk.Msg = MsgChannelOpenTry{}

type MsgChannelOpenTry struct {
	PortID              string            `json:"port_id"`
	ChannelID           string            `json:"channel_id"`
	Channel             Channel           `json:"channel"`
	CounterpartyVersion string            `json:"counterparty_version"`
	ProofInit           commitment.ProofI `json:"proof_init"`
	ProofHeight         uint64            `json:"proof_height"`
	Signer              sdk.AccAddress    `json:"signer"`
}

// NewMsgChannelOpenTry creates a new MsgChannelOpenTry instance
func NewMsgChannelOpenTry(
	portID, channelID, version string, channelOrder Order, connectionHops []string,
	counterpartyPortID, counterpartyChannelID, counterpartyVersion string,
	proofInit commitment.ProofI, proofHeight uint64, signer sdk.AccAddress,
) MsgChannelOpenTry {
	counterparty := NewCounterparty(counterpartyPortID, counterpartyChannelID)
	channel := NewChannel(INIT, channelOrder, counterparty, connectionHops, version)
	return MsgChannelOpenTry{
		PortID:              portID,
		ChannelID:           channelID,
		Channel:             channel,
		CounterpartyVersion: counterpartyVersion,
		ProofInit:           proofInit,
		ProofHeight:         proofHeight,
		Signer:              signer,
	}
}

// Route implements sdk.Msg
func (msg MsgChannelOpenTry) Route() string {
	return ibctypes.RouterKey
}

// Type implements sdk.Msg
func (msg MsgChannelOpenTry) Type() string {
	return "channel_open_try"
}

// ValidateBasic implements sdk.Msg
func (msg MsgChannelOpenTry) ValidateBasic() sdk.Error {
	if err := host.DefaultConnectionIdentifierValidator(msg.PortID); err != nil {
		return sdk.ConvertError(sdkerrors.Wrap(err, "invalid port ID"))
	}
	if err := host.DefaultClientIdentifierValidator(msg.ChannelID); err != nil {
		return sdk.ConvertError(sdkerrors.Wrap(err, "invalid channel ID"))
	}
	if strings.TrimSpace(msg.CounterpartyVersion) == "" {
		return sdk.ConvertError(ErrInvalidCounterpartyChannel(DefaultCodespace, "counterparty version cannot be blank"))
	}
	if msg.ProofInit == nil {
		return sdk.ConvertError(ibctypes.ErrInvalidProof(DefaultCodespace, "cannot submit an empty proof"))
	}
	if err := msg.ProofInit.ValidateBasic(); err != nil {
		return sdk.ConvertError(ibctypes.ErrInvalidProof(DefaultCodespace, err.Error()))
	}
	if msg.ProofHeight == 0 {
		return sdk.ConvertError(ibctypes.ErrInvalidProof(DefaultCodespace, "proof height must be > 0"))
	}
	// Signer can be empty
	return sdk.ConvertError(msg.Channel.ValidateBasic())
}

// GetSignBytes implements sdk.Msg
func (msg MsgChannelOpenTry) GetSignBytes() []byte {
	return sdk.MustSortJSON(SubModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgChannelOpenTry) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

var _ sdk.Msg = MsgChannelOpenAck{}

type MsgChannelOpenAck struct {
	PortID              string            `json:"port_id"`
	ChannelID           string            `json:"channel_id"`
	CounterpartyVersion string            `json:"counterparty_version"`
	ProofTry            commitment.ProofI `json:"proof_try"`
	ProofHeight         uint64            `json:"proof_height"`
	Signer              sdk.AccAddress    `json:"signer"`
}

// NewMsgChannelOpenAck creates a new MsgChannelOpenAck instance
func NewMsgChannelOpenAck(
	portID, channelID string, cpv string, proofTry commitment.ProofI, proofHeight uint64,
	signer sdk.AccAddress,
) MsgChannelOpenAck {
	return MsgChannelOpenAck{
		PortID:              portID,
		ChannelID:           channelID,
		CounterpartyVersion: cpv,
		ProofTry:            proofTry,
		ProofHeight:         proofHeight,
		Signer:              signer,
	}
}

// Route implements sdk.Msg
func (msg MsgChannelOpenAck) Route() string {
	return ibctypes.RouterKey
}

// Type implements sdk.Msg
func (msg MsgChannelOpenAck) Type() string {
	return "channel_open_ack"
}

// ValidateBasic implements sdk.Msg
func (msg MsgChannelOpenAck) ValidateBasic() sdk.Error {
	if err := host.DefaultConnectionIdentifierValidator(msg.PortID); err != nil {
		return sdk.ConvertError(sdkerrors.Wrap(err, "invalid port ID"))
	}
	if err := host.DefaultClientIdentifierValidator(msg.ChannelID); err != nil {
		return sdk.ConvertError(sdkerrors.Wrap(err, "invalid channel ID"))
	}
	if strings.TrimSpace(msg.CounterpartyVersion) == "" {
		return sdk.ConvertError(ErrInvalidCounterpartyChannel(DefaultCodespace, "counterparty version cannot be blank"))
	}
	if msg.ProofTry == nil {
		return sdk.ConvertError(ibctypes.ErrInvalidProof(DefaultCodespace, "cannot submit an empty proof"))
	}
	if err := msg.ProofTry.ValidateBasic(); err != nil {
		return sdk.ConvertError(ibctypes.ErrInvalidProof(DefaultCodespace, err.Error()))
	}
	if msg.ProofHeight == 0 {
		return sdk.ConvertError(ibctypes.ErrInvalidProof(DefaultCodespace, "proof height must be > 0"))
	}
	// Signer can be empty
	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgChannelOpenAck) GetSignBytes() []byte {
	return sdk.MustSortJSON(SubModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgChannelOpenAck) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

var _ sdk.Msg = MsgChannelOpenConfirm{}

type MsgChannelOpenConfirm struct {
	PortID      string            `json:"port_id"`
	ChannelID   string            `json:"channel_id"`
	ProofAck    commitment.ProofI `json:"proof_ack"`
	ProofHeight uint64            `json:"proof_height"`
	Signer      sdk.AccAddress    `json:"signer"`
}

// NewMsgChannelOpenConfirm creates a new MsgChannelOpenConfirm instance
func NewMsgChannelOpenConfirm(
	portID, channelID string, proofAck commitment.ProofI, proofHeight uint64,
	signer sdk.AccAddress,
) MsgChannelOpenConfirm {
	return MsgChannelOpenConfirm{
		PortID:      portID,
		ChannelID:   channelID,
		ProofAck:    proofAck,
		ProofHeight: proofHeight,
		Signer:      signer,
	}
}

// Route implements sdk.Msg
func (msg MsgChannelOpenConfirm) Route() string {
	return ibctypes.RouterKey
}

// Type implements sdk.Msg
func (msg MsgChannelOpenConfirm) Type() string {
	return "channel_open_confirm"
}

// ValidateBasic implements sdk.Msg
func (msg MsgChannelOpenConfirm) ValidateBasic() sdk.Error {
	if err := host.DefaultConnectionIdentifierValidator(msg.PortID); err != nil {
		return sdk.ConvertError(sdkerrors.Wrap(err, "invalid port ID"))
	}
	if err := host.DefaultClientIdentifierValidator(msg.ChannelID); err != nil {
		return sdk.ConvertError(sdkerrors.Wrap(err, "invalid channel ID"))
	}
	if msg.ProofAck == nil {
		return sdk.ConvertError(ibctypes.ErrInvalidProof(DefaultCodespace, "cannot submit an empty proof"))
	}
	if err := msg.ProofAck.ValidateBasic(); err != nil {
		return sdk.ConvertError(ibctypes.ErrInvalidProof(DefaultCodespace, err.Error()))
	}
	if msg.ProofHeight == 0 {
		return sdk.ConvertError(ibctypes.ErrInvalidProof(DefaultCodespace, "proof height must be > 0"))
	}
	// Signer can be empty
	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgChannelOpenConfirm) GetSignBytes() []byte {
	return sdk.MustSortJSON(SubModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgChannelOpenConfirm) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

var _ sdk.Msg = MsgChannelCloseInit{}

type MsgChannelCloseInit struct {
	PortID    string         `json:"port_id"`
	ChannelID string         `json:"channel_id"`
	Signer    sdk.AccAddress `json:"signer"`
}

// NewMsgChannelCloseInit creates a new MsgChannelCloseInit instance
func NewMsgChannelCloseInit(portID string, channelID string, signer sdk.AccAddress) MsgChannelCloseInit {
	return MsgChannelCloseInit{
		PortID:    portID,
		ChannelID: channelID,
		Signer:    signer,
	}
}

// Route implements sdk.Msg
func (msg MsgChannelCloseInit) Route() string {
	return ibctypes.RouterKey
}

// Type implements sdk.Msg
func (msg MsgChannelCloseInit) Type() string {
	return "channel_close_init"
}

// ValidateBasic implements sdk.Msg
func (msg MsgChannelCloseInit) ValidateBasic() sdk.Error {
	if err := host.DefaultConnectionIdentifierValidator(msg.PortID); err != nil {
		return sdk.ConvertError(sdkerrors.Wrap(err, "invalid port ID"))
	}
	if err := host.DefaultClientIdentifierValidator(msg.ChannelID); err != nil {
		return sdk.ConvertError(sdkerrors.Wrap(err, "invalid channel ID"))
	}
	// Signer can be empty
	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgChannelCloseInit) GetSignBytes() []byte {
	return sdk.MustSortJSON(SubModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgChannelCloseInit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

var _ sdk.Msg = MsgChannelCloseConfirm{}

type MsgChannelCloseConfirm struct {
	PortID      string            `json:"port_id"`
	ChannelID   string            `json:"channel_id"`
	ProofInit   commitment.ProofI `json:"proof_init"`
	ProofHeight uint64            `json:"proof_height"`
	Signer      sdk.AccAddress    `json:"signer"`
}

// NewMsgChannelCloseConfirm creates a new MsgChannelCloseConfirm instance
func NewMsgChannelCloseConfirm(
	portID, channelID string, proofInit commitment.ProofI, proofHeight uint64,
	signer sdk.AccAddress,
) MsgChannelCloseConfirm {
	return MsgChannelCloseConfirm{
		PortID:      portID,
		ChannelID:   channelID,
		ProofInit:   proofInit,
		ProofHeight: proofHeight,
		Signer:      signer,
	}
}

// Route implements sdk.Msg
func (msg MsgChannelCloseConfirm) Route() string {
	return ibctypes.RouterKey
}

// Type implements sdk.Msg
func (msg MsgChannelCloseConfirm) Type() string {
	return "channel_close_confirm"
}

// ValidateBasic implements sdk.Msg
func (msg MsgChannelCloseConfirm) ValidateBasic() sdk.Error {
	if err := host.DefaultPortIdentifierValidator(msg.PortID); err != nil {
		return sdk.ConvertError(sdkerrors.Wrap(err, "invalid port ID"))
	}
	if err := host.DefaultChannelIdentifierValidator(msg.ChannelID); err != nil {
		return sdk.ConvertError(sdkerrors.Wrap(err, "invalid channel ID"))
	}
	if msg.ProofInit == nil {
		return sdk.ConvertError(ibctypes.ErrInvalidProof(DefaultCodespace, "cannot submit an empty proof"))
	}
	if err := msg.ProofInit.ValidateBasic(); err != nil {
		return sdk.ConvertError(ibctypes.ErrInvalidProof(DefaultCodespace, err.Error()))
	}
	if msg.ProofHeight == 0 {
		return sdk.ConvertError(ibctypes.ErrInvalidProof(DefaultCodespace, "proof height must be > 0"))
	}
	// Signer can be empty
	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgChannelCloseConfirm) GetSignBytes() []byte {
	return sdk.MustSortJSON(SubModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgChannelCloseConfirm) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// MsgPacket receives incoming IBC packet
type MsgPacket struct {
	Packet `json:"packet" yaml:"packet"`

	Proof       commitment.ProofI `json:"proof" yaml:"proof"`
	ProofHeight uint64            `json:"proof_height" yaml:"proof_height"`

	Signer sdk.AccAddress `json:"signer" yaml:"signer"`
}

var _ sdk.Msg = MsgPacket{}

// NewMsgPacket constructs new MsgPacket
func NewMsgPacket(packet Packet, proof commitment.ProofI, proofHeight uint64, signer sdk.AccAddress) MsgPacket {
	return MsgPacket{
		Packet: packet,

		Proof:       proof,
		ProofHeight: proofHeight,

		Signer: signer,
	}
}

// Implements sdk.Msg
func (msg MsgPacket) Route() string {
	return msg.DestinationPort
}

// Implements sdk.Msg
func (msg MsgPacket) ValidateBasic() sdk.Error {
	if msg.Proof == nil {
		return sdk.ConvertError(ibctypes.ErrInvalidProof(DefaultCodespace, "cannot submit an empty proof"))
	}
	if err := msg.Proof.ValidateBasic(); err != nil {
		return sdk.ConvertError(ibctypes.ErrInvalidProof(DefaultCodespace, err.Error()))
	}
	if msg.ProofHeight == 0 {
		return sdk.ConvertError(ibctypes.ErrInvalidProof(DefaultCodespace, "proof height must be > 0"))
	}
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress("missing signer address")
	}

	return sdk.ConvertError(msg.PacketDataI.ValidateBasic())
}

// Implements sdk.Msg
func (msg MsgPacket) GetSignBytes() []byte {
	return sdk.MustSortJSON(SubModuleCdc.MustMarshalJSON(msg))
}

// Implements sdk.Msg
func (msg MsgPacket) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

var _ sdk.Msg = MsgTimeout{}

// MsgTimeout receives timeouted packet
type MsgTimeout struct {
	Packet           `json:"packet" yaml:"packet"`
	NextSequenceRecv uint64 `json:"next_sequence_recv,omitempty" yaml:"next_sequence_recv,omitempty"`

	Proof       commitment.ProofI `json:"proof" yaml:"proof"`
	ProofHeight uint64            `json:"proof_height" yaml:"proof_height"`

	Signer sdk.AccAddress `json:"signer" yaml:"signer"`
}

// NewMsgTimeout constructs new MsgTimeout
func NewMsgTimeout(packet Packet, nextSequenceRecv uint64, proof commitment.ProofI, proofHeight uint64, signer sdk.AccAddress) MsgTimeout {
	return MsgTimeout{
		Packet:           packet,
		NextSequenceRecv: nextSequenceRecv,
		Proof:            proof,
		ProofHeight:      proofHeight,
		Signer:           signer,
	}
}

// Implements sdk.Msg
func (msg MsgTimeout) Route() string {
	return msg.SourcePort
}

// Implements sdk.Msg
func (msg MsgTimeout) ValidateBasic() sdk.Error {
	if msg.Proof == nil {
		return sdk.ConvertError(ibctypes.ErrInvalidProof(DefaultCodespace, "cannot submit an empty proof"))
	}
	if err := msg.Proof.ValidateBasic(); err != nil {
		return sdk.ConvertError(ibctypes.ErrInvalidProof(DefaultCodespace, err.Error()))
	}
	if msg.ProofHeight == 0 {
		return sdk.ConvertError(ibctypes.ErrInvalidProof(DefaultCodespace, "proof height must be > 0"))
	}
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress("missing signer address")
	}

	return sdk.ConvertError(msg.PacketDataI.ValidateBasic())
}

// Implements sdk.Msg
func (msg MsgTimeout) GetSignBytes() []byte {
	return sdk.MustSortJSON(SubModuleCdc.MustMarshalJSON(msg))
}

// Implements sdk.Msg
func (msg MsgTimeout) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

var _ sdk.Msg = MsgAcknowledgement{}

// MsgAcknowledgement receives incoming IBC acknowledgement
type MsgAcknowledgement struct {
	Packet          `json:"packet" yaml:"packet"`
	Acknowledgement PacketDataI `json:"acknowledgement" yaml:"acknowledgement"`

	Proof       commitment.ProofI `json:"proof" yaml:"proof"`
	ProofHeight uint64            `json:"proof_height" yaml:"proof_height"`

	Signer sdk.AccAddress `json:"signer" yaml:"signer"`
}

// NewMsgAcknowledgement constructs a new MsgAcknowledgement
func NewMsgAcknowledgement(packet Packet, ack PacketDataI, proof commitment.ProofI, proofHeight uint64, signer sdk.AccAddress) MsgAcknowledgement {
	return MsgAcknowledgement{
		Packet:          packet,
		Acknowledgement: ack,
		Proof:           proof,
		ProofHeight:     proofHeight,
		Signer:          signer,
	}
}

// Implements sdk.Msg
func (msg MsgAcknowledgement) Route() string {
	return msg.SourcePort
}

// Implements sdk.Msg
func (msg MsgAcknowledgement) ValidateBasic() sdk.Error {
	if msg.Proof == nil {
		return sdk.ConvertError(ibctypes.ErrInvalidProof(DefaultCodespace, "cannot submit an empty proof"))
	}
	if msg.ProofHeight == 0 {
		return sdk.ConvertError(ibctypes.ErrInvalidProof(DefaultCodespace, "proof height must be > 0"))
	}

	if err := msg.Packet.ValidateBasic(); err != nil {
		return sdk.ConvertError(err)
	}

	if err := msg.Acknowledgement.ValidateBasic(); err != nil {
		return sdk.ConvertError(err)
	}

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress("missing signer address")
	}

	return nil
}

// Implements sdk.Msg
func (msg MsgAcknowledgement) GetSignBytes() []byte {
	return sdk.MustSortJSON(SubModuleCdc.MustMarshalJSON(msg))
}

// Implements sdk.Msg
func (msg MsgAcknowledgement) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
