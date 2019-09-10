package actions

// Ping returns a status of a ping request. Requires a FQDN (e.g. http:// or https://)
func Ping(FQDN string) {
	GetDestinationStatus("http://<url/to/ping>/")
}
