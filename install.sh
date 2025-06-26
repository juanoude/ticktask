go build
sudo mv -f ticktask /usr/local/bin/ticktask

if [ -d $HOME/.ticktask/music ]; then
    sudo cp -r music/focus $HOME/.ticktask/music/
    sudo cp -r music/idle $HOME/.ticktask/music/
    sudo cp -r music/generic $HOME/.ticktask/music/
else
    mkdir $HOME/.ticktask/music
    sudo cp -r music/focus $HOME/.ticktask/music/
    sudo cp -r music/idle $HOME/.ticktask/music/
    sudo cp -r music/generic $HOME/.ticktask/music/
fi
