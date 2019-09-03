/**

MIT License

Copyright (c) 2019 Job Talle

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

const BodyShape = function() {
    const length = BodyShape.LENGTH_MIN + (BodyShape.LENGTH_MAX - BodyShape.LENGTH_MIN) * Math.random();
    const widths = [];
    const legAngle = BodyShape.LEG_ANGLE_MIN + (BodyShape.LEG_ANGLE_MAX - BodyShape.LEG_ANGLE_MIN) * Math.random();
    const legFactor = BodyShape.LEG_FACTOR_MIN + (BodyShape.LEG_FACTOR_MAX - BodyShape.LEG_FACTOR_MIN) * Math.random();
    const thicknessMultiplier = BodyShape.THICKNESS_MULTIPLIER_MIN + (1 - BodyShape.THICKNESS_MULTIPLIER_MIN) * Math.random();
    const fill = BodyShape.COLORS[Math.floor(Math.random() * BodyShape.COLORS.length)];
    let thickness = 0;

    const getWidthMultiplier = i => {
        return Math.sin((i / (widths.length - 1)) * Math.PI * BodyShape.SINE_WAVE_PORTION);
    };

    this.getLength = () => (widths.length + 1) * BodyShape.LENGTH_SEGMENT;

    this.getThickness = () => thickness;

    this.getLegAngle = () => legAngle;

    this.getLegLength = () => Math.max(BodyShape.LEG_LENGTH_MIN, thickness * BodyShape.LEG_SCALE * legFactor);

    this.draw = context => {
        const right = this.getLength() * 0.5;

        context.fillStyle = fill;
        context.strokeStyle = "black";
        context.beginPath();
        context.moveTo(right, 0);

        for (let i = 0; i < widths.length; ++i)
            context.lineTo(
                right - i * BodyShape.LENGTH_SEGMENT,
                -widths[i] * getWidthMultiplier(i));

        for (let i = widths.length; i-- > 0;)
            context.lineTo(
                right - i * BodyShape.LENGTH_SEGMENT,
                widths[i] * getWidthMultiplier(i));

        context.closePath();
        context.fill();
        context.stroke();
    };

    while (this.getLength() < length) {
        const newThickness = (BodyShape.WIDTH_MIN + (BodyShape.WIDTH_MAX - BodyShape.WIDTH_MIN) * Math.random()) * thicknessMultiplier;

        if (thickness < newThickness)
            thickness = newThickness;

        widths.push(newThickness);
    }
};

BodyShape.WIDTH_MIN = 12;
BodyShape.WIDTH_MAX = 32;
BodyShape.THICKNESS_MULTIPLIER_MIN = 0.3;
BodyShape.LENGTH_MIN = 18;
BodyShape.LENGTH_MAX = 80;
BodyShape.LENGTH_SEGMENT = 6;
BodyShape.LEG_ANGLE_MIN = 0.3;
BodyShape.LEG_ANGLE_MAX = 0.9;
BodyShape.SINE_WAVE_PORTION = 0.9;
BodyShape.LEG_SCALE = 1.8;
BodyShape.LEG_LENGTH_MIN = 18;
BodyShape.LEG_FACTOR_MIN = 0.5;
BodyShape.LEG_FACTOR_MAX = 1;
BodyShape.COLORS = [
    "rgba(107,142,35,0.7)",
    "rgba(178,34,34,0.7)",
    "rgba(222,184,135,0.7)",
    "rgba(95,158,160,0.7)"
];
