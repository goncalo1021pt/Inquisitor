package inquisitor

func exec() error {
	go signalHandler()

	for {
		select {
		case <-shutdownSignal:
			return nil
		default:
			if err := arp(); err != nil {
				return err
			}
			if err := ftp(); err != nil {
				return err
			}
		}
	}
}
