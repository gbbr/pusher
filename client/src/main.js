var canvas = document.getElementById("surface"),
	canvas2 = document.getElementById("incoming"),

	log = function (elementId, value) {
		document.getElementById(elementId).innerHTML = value;
	},

	getRandomColor = function () {
		var letters = '0123456789ABCDEF'.split('');
		var color = '#';
		for (var i = 0; i < 6; i++ ) {
			color += letters[Math.floor(Math.random() * 16)];
		}
		return color;
	},

	drawCircle = function (canvas, circle) {
		var context = canvas.getContext("2d");

		// context.clearRect(0, 0, canvas.width, canvas.height);

		context.beginPath();
			context.arc(circle.x, circle.y, circle.r, 0, 2 * Math.PI, false);
			context.fillStyle = circle.c;
			context.fill();
			context.lineWidth = 5;
			context.strokeStyle = "#003300";
		context.stroke();
	},

	color = getRandomColor();


var socket = new WebSocket("ws://mecca.local:888/pipe");

socket.onopen = function (event) {
	new PressureDetector(canvas).attachCallback(function (position, pressure) {
		if (pressure === 0) color = getRandomColor();

		var circle = {
			x: position.x,
			y: position.y,
			r: pressure,
			c: color
		};

		socket.send(JSON.stringify(circle));

		log("valueX", position.x);
		log("valueY", position.y);
		log("strength", pressure);
	});

	socket.onmessage = function (msg) {
		var circleData = JSON.parse(msg.data);
		drawCircle(canvas, circleData);
	}
};

socket.onclose = function () {
	console.log("Socket CLOSED");
}
