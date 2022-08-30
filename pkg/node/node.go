package node

type Slot struct {
	Period int `json:"period"`
	Thread int `json:"thread"`
}

func NextSlot(c *Client) (uint64, error) {
	status, err := Status(c)
	if err != nil {
		return 0, err
	}

	return uint64(status.NextSlot.Period), nil
}
