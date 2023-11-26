const video = document.getElementsByTagName('video')[0];
const getDisplay = document.getElementById('get-display');
let videoSettings;
getDisplay.onclick = function () {
    navigator.mediaDevices.getDisplayMedia({
        video: true,
        audio: false
    }).then(stream => {
        videoSettings = stream.getVideoTracks()[0].getSettings();
        video.srcObject = stream;
        window.stream = stream;
    }).catch(err => {
        console.error(err);
    });
}

const record = document.getElementById('record');
const storeRecord = document.getElementById('stop-record');
let recordedBlobs = [];

function onStartRecord() {
    record.hidden = true;
    storeRecord.hidden = false;
}

function onStopRecord() {
    storeRecord.hidden = true;
    record.hidden = false;
}

record.onclick = function () {
    onStartRecord();
    const options = {mimeType: 'video/webm'};
    const mediaRecorder = new MediaRecorder(window.stream, options);
    mediaRecorder.start(10);
    mediaRecorder.ondataavailable = function (event) {
        if (event.data && event.data.size > 0) {
            recordedBlobs.push(event.data);
        }
    }
}

let corpInfo = {};

function download() {
    // get format
    let formatElement = document.getElementsByName('output-format');
    let format = "mp4"
    for (let i = 0; i < formatElement.length; i++) {
        if (formatElement[i].checked) {
            corpInfo.format = formatElement[i].value;
            format = formatElement[i].value;
            break;
        }
    }

    const blob = new Blob(recordedBlobs, {type: 'video/webm'});
    let formData = new FormData();
    formData.append('file', blob);
    formData.append('corpInfo', JSON.stringify(corpInfo));

    fetch('/api/download', {
        method: 'POST',
        body: formData
    }).then(res => {
        if (res.status === 500) {
            res.text().then(text => {
                alert('裁剪失败: ' + text);
            })
            return;
        }
        res.blob().then(blob => {
            const a = document.createElement('a');
            a.style.display = 'none';
            a.href = window.URL.createObjectURL(blob);
            a.download = `download.${format}`;
            document.body.appendChild(a);
            a.click();
        })
    }).catch(err => {
        console.log(err)
    }).finally(() => {
        corpInfo = {};
        recordedBlobs = [];
    })
}

storeRecord.onclick = () => {
    onStopRecord();
    let tracks = video.srcObject.getTracks();
    tracks.forEach(track => {
        track.stop();
    });
    video.srcObject = null;
    download();
}

// 获取目标元素和选框元素
const target = document.getElementById('video');
const selectionBox = document.getElementById('selection-box');
let startX, startY, endX, endY
target.addEventListener('mousedown', (e) => {
    startX = e.clientX;
    startY = e.clientY;
    selectionBox.style.left = startX + 'px';
    selectionBox.style.top = startY + 'px';
    selectionBox.style.display = 'block';
});
// 监听鼠标移动事件
target.addEventListener('mousemove', (e) => {
    if (startX !== null && startY !== null) {
        endX = e.clientX;
        endY = e.clientY;
        selectionBox.style.width = Math.abs(e.clientX - startX) + 'px';
        selectionBox.style.height = Math.abs(e.clientY - startY) + 'px';
        selectionBox.style.left = e.clientX < startX ? e.clientX : startX + 'px';
        selectionBox.style.top = e.clientY < startY ? e.clientY : startY + 'px';
    }
});

// 监听鼠标松开事件
target.addEventListener('mouseup', () => {
    let widthRatio = video.clientWidth / videoSettings.width;
    let heightRatio = video.clientHeight / videoSettings.height;
    corpInfo = {
        "x": Math.min(startX, endX) / widthRatio, "y": Math.min(startY, endY) / heightRatio,
        "width": selectionBox.offsetWidth / widthRatio, "height": selectionBox.offsetHeight / heightRatio
    };
    startX = null;
    startY = null;
    selectionBox.style.width = '0';
    selectionBox.style.height = '0';
    selectionBox.style.left = '0';
    selectionBox.style.top = '0';
    selectionBox.style.display = 'none';
});
