/**
 * @description Binds mouse events to a DOM element, detecting
 * position and pressure
 * @param el
 * @constructor
 */
function PressureDetector(el) {
	this.pressure = 0;

	el.addEventListener("mousedown", this.applyPressure.bind(this));
	el.addEventListener("mouseup", this.releasePressure.bind(this));
	el.addEventListener("mouseleave", this.releasePressure.bind(this));

	el.addEventListener("touchstart", this.applyPressure.bind(this));
	el.addEventListener("touchend", this.releasePressure.bind(this));
}

PressureDetector.prototype = {
	MAX_PRESSURE      : 10,
	INCREASE_SPEED_MS : 10,

	// Canvas object
	pressure: null,
	position: null,

	// Invoked when pressure changes
	_listener: null,
	_interval: null,

	applyPressure: function (event) {
		this.position = {
			x: event.layerX,
			y: event.layerY
		};

		clearInterval(this._interval);

		this._interval = setInterval(function() {
			if (this.pressure >= this.MAX_PRESSURE) {
				clearInterval(this._interval);
			} else {
				this.setPressure(this.pressure + 1);
			}
		}.bind(this), this.INCREASE_SPEED_MS);
	},

	releasePressure: function () {
		clearInterval(this._interval);

		this._interval = setInterval(function() {
			if (this.pressure == 0) {
				clearInterval(this._interval);
			} else {
				this.setPressure(this.pressure - 1);
			}
		}.bind(this), this.INCREASE_SPEED_MS);
	},

	setPressure: function (v) {
		this.pressure = v;

		if (this._listener) {
			this._listener(this.position, this.pressure);
		}
	},

	attachCallback: function (fn) {
		this._listener = fn;
	}
};