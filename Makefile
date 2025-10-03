run:
	touch ~/.config/fyne/dev.jkulzer.transit-tool/realtimeGtfs.bin
	rm ~/.config/fyne/dev.jkulzer.transit-tool/realtimeGtfs.bin
	go run .
android-install:
	fyne package -os android -app-id com.example.myapp
	# adb install osm_test.apk
