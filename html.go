package main

func html(serverHost string, serverPort string) string {
	return `<!DOCTYPE HTML>
	<html lang="en">
	<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	<title>Huge File Reader</title>
	</head>
	<body style="margin:0;font-family:sans-serif;background:#4c8;">
		<div id="progress">Drop a file...</div>
		<script>` + jsmd5() + `</script>
		<script>
	
	(function() {
	var chunkSize = 1024 * 1024 * 4
	var socket = new WebSocket('ws://` + serverHost + ":" + serverPort + `/ws')
	socket.onmessage = function() {
		console.log('Received ready event')
	}
	var closeSocket = function() {
		if (socket.readyState !== 1) {
			console.log(socket.readyState)
			socket.close()
			setTimeout(function() {
				socket = new WebSocket('ws://` + serverHost + ":" + serverPort + `/ws')
				if (typeof socket.onclose === 'undefined') socket.onclose = closeSocket
				socket.onmessage = function() {
					console.log('Received ready event')
				}
			}, 3000)
			}
	}
	if (typeof socket.onclose === 'undefined') socket.onclose = closeSocket
	var progress = document.getElementById('progress')
	progress.style.background = '#396'
	progress.style.color = '#fff'
	progress.style.fontWeight = 'bold'
	progress.style.padding = '5px'
	progress.style.display = 'block'
	progress.style.height = '20px'
	progress.style.whiteSpace = 'nowrap'
	function parseFile(file, options) {
		var fileSize = file.size
		var offset = 0
		var readBlock = null
		var chunkReadCallback = function(data) {
			console.log(data)
			socket.send(data)
		}
		var chunkErrorCallback = function(err) {
			console.log('ERROR', err)
		}
		var result = function(msg, count) {
			console.log(msg + ' ' + count)
		}
		var onLoadHandler = function(evt) {
			if (evt.target.error == null) {
				offset += evt.loaded
				chunkReadCallback(evt.target.result)
			} else {
				chunkErrorCallback(evt.target.error)
				return
			}
			var percentage = Math.round((offset / fileSize) * 100)
			progress.innerHTML = file.name + ' (MD5=' + socket.uploadChecksum + ') ' +
				' &nbsp; ' + percentage + '% uploaded'
			progress.style.width = percentage + '%'
			if (offset === fileSize) {
				result('Success!', offset)
				socket.send('ready:' +  offset)
				return
			} else if (offset > fileSize) {
				result('Fail!', offset)
				return
			}
			readBlock(offset, chunkSize, file)
		}
		readBlock = function(_offset, length, _file) {
			var r = new FileReader()
			var blob = _file.slice(_offset, length + _offset);
			r.onload = onLoadHandler
			r.readAsArrayBuffer(blob)
		}
		readBlock(offset, chunkSize, file)
	}
	
	function getFileChecksum(file, options) {
		var fileSize = file.size
		var spark = new SparkMD5.ArrayBuffer()
		var offset = 0
		var readBlock = null
		var chunkReadCallback = function(data) {
			spark.append(data)
		}
		var chunkErrorCallback = function(err) {
			console.log('ERROR', err)
		}
		var onLoadHandler = function(evt) {
			if (evt.target.error == null) {
				offset += evt.loaded
				chunkReadCallback(evt.target.result)
			} else {
				chunkErrorCallback(evt.target.error)
				return
			}
			var percentage = Math.round((offset / fileSize) * 100)
			progress.innerHTML = 'Calculating MD5 for ' + file.name + ' &nbsp; ' + percentage + '%'
			progress.style.width = percentage + '%'
			if (offset === fileSize) {
				socket.uploadChecksum = spark.end()
				parseFile(file)
				return
			} else if (offset > fileSize) {
				return
			}
			readBlock(offset, chunkSize, file)
		}
		readBlock = function(_offset, length, _file) {
			var r = new FileReader()
			var blob = _file.slice(_offset, length + _offset);
			r.onload = onLoadHandler
			r.readAsArrayBuffer(blob)
		}
		readBlock(offset, chunkSize, file)
	}
	
	window.ondragover = function() { return false }
	window.ondrop = function(e) { 
		if (e.dataTransfer.files.length > 0) {
			socket.send('upload:' +  e.dataTransfer.files[0].name)
			socket.uploadChecksum = () => {}
			getFileChecksum(e.dataTransfer.files[0])
		}
		return false 
	}
	})()
		</script>
	</body>
	</html>
	`
}
