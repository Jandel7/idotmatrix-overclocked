# iDotMatrix Controller — Quick Reference

## Device Info
- **BLE Address (Mac):** `99961A10-7BC3-D4BA-0C65-41334CA19A46`
- **BLE Address (Pi):** `56:7B:B9:3E:57:3A`
- **Pi IP:** `192.168.1.113`
- **Pi User:** `pi`

---

## Running from the Raspberry Pi

### SSH into the Pi
```
ssh pi@192.168.1.113
```

### Launch the menu
```
idm
```

### Or run individual commands directly
```
sudo /home/pi/idotmatrix-overclocked/idm-cli --target 56:7B:B9:3E:57:3A fire
sudo /home/pi/idotmatrix-overclocked/idm-cli --target 56:7B:B9:3E:57:3A snake
sudo /home/pi/idotmatrix-overclocked/idm-cli --target 56:7B:B9:3E:57:3A tetris
sudo /home/pi/idotmatrix-overclocked/idm-cli --target 56:7B:B9:3E:57:3A clock
sudo /home/pi/idotmatrix-overclocked/idm-cli --target 56:7B:B9:3E:57:3A demo
sudo /home/pi/idotmatrix-overclocked/idm-cli --target 56:7B:B9:3E:57:3A emoji --name rocket
sudo /home/pi/idotmatrix-overclocked/idm-cli --target 56:7B:B9:3E:57:3A text --text "HELLO" --color white --animation fireworks
```

---

## Running from your Mac

### Navigate to the project folder
```
cd ~/iDotMatrix/idotmatrix-overclocked
```

### Run individual commands
```
./idm-cli fire --target 99961A10-7BC3-D4BA-0C65-41334CA19A46
./idm-cli snake --target 99961A10-7BC3-D4BA-0C65-41334CA19A46
./idm-cli tetris --target 99961A10-7BC3-D4BA-0C65-41334CA19A46
./idm-cli clock --target 99961A10-7BC3-D4BA-0C65-41334CA19A46
./idm-cli demo --target 99961A10-7BC3-D4BA-0C65-41334CA19A46
./idm-cli emoji --name rocket --target 99961A10-7BC3-D4BA-0C65-41334CA19A46
./idm-cli text --text "HELLO" --color white --animation fireworks --target 99961A10-7BC3-D4BA-0C65-41334CA19A46
```

### Snake and Tetris require sudo on Mac for keyboard input
```
sudo ./idm-cli snake --target 99961A10-7BC3-D4BA-0C65-41334CA19A46
sudo ./idm-cli tetris --target 99961A10-7BC3-D4BA-0C65-41334CA19A46
```

---

## Weather Station

### Start manually
```
ssh pi@192.168.1.113
cd /home/pi/idotmatrix-weather
source .venv/bin/activate
python3 weather_display.py
```

### Start/stop the auto-updating service
```
sudo systemctl start weather.service
sudo systemctl stop weather.service
sudo systemctl status weather.service
```

### Enable/disable on boot
```
sudo systemctl enable weather.service
sudo systemctl disable weather.service
```

---

## Rebuilding after code changes

### On the Pi
```
ssh pi@192.168.1.113
cd /home/pi/idotmatrix-overclocked
make build
```

### On the Mac
```
cd ~/iDotMatrix/idotmatrix-overclocked
make build
```

---

## Controls

### Snake
| Key | Action |
|-----|--------|
| W / ↑ | Move Up |
| S / ↓ | Move Down |
| A / ← | Move Left |
| D / → | Move Right |
| Q | Quit |

### Tetris
| Key | Action |
|-----|--------|
| A / ← | Move Left |
| D / → | Move Right |
| W / ↑ | Rotate |
| S / ↓ | Soft Drop |
| Space | Hard Drop |
| Q | Quit |

---

## TODO — Future Ideas
- [ ] Web controller accessible from iPhone browser
- [ ] D-pad on screen for Snake and Tetris
- [ ] Launch games and animations without SSH
