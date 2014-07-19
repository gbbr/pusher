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
	MAX_PRESSURE      : 40,
	INCREASE_SPEED_MS : 10,

	// Canvas object
	pressure: null,
	position: null,

	// Invoked when pressure changes
	_listener: null,

	applyPressure: function (event) {
		this.position = {
			x: event.layerX,
			y: event.layerY
		};

		this.setPressure(1);
	},

	releasePressure: function () {
		this.setPressure(0);
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