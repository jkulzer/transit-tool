remove-realtime:
	touch ~/.config/fyne/dev.jkulzer.transit-tool/realtimeGtfs.bin
	rm ~/.config/fyne/dev.jkulzer.transit-tool/realtimeGtfs.bin
run: remove-realtime
	go run .
debug: remove-realtime
	go run . debug
android-install:
	fyne package -os android -app-id com.example.myapp
	# adb install osm_test.apk
