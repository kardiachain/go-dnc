package tx_pool

func StringifyBlacklist() string {
	addrs := ""
	for addr, _ := range Blacklisted {
		addrs = addrs + addr + ","
	}
	return addrs
}
