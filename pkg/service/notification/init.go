package notification

type NotiService interface {
	SendNotiToDevice(message string) error
}
