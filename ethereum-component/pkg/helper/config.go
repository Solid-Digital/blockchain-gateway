package helper

import "fmt"

func (h *Helper) FillInTriggerConfig(configFmt, contractAddress, ABI string) []byte {
	return []byte(fmt.Sprintf(configFmt, contractAddress, ABI, contractAddress))
}

func (h *Helper) FillInActionConfig(configFmt, contractAddress, ABI string) []byte {
	return []byte(fmt.Sprintf(configFmt, contractAddress, ABI))
}
