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

const Leg = function(x, y, bugDirection, direction, length, initialProgress) {
    let counterpart = null;
    let onGround = true;
    let footX = x + Math.cos(bugDirection + direction) * length;
    let footY = y + Math.sin(bugDirection + direction) * length;

    const applyInitialProgress = () => {
        const dist = Math.cos(direction) * length * -2 * initialProgress;

        footX += Math.cos(bugDirection) * dist;
        footY += Math.sin(bugDirection) * dist;
    };

    this.setCounterpart = leg => counterpart = leg;

    this.stepDown = () => {
        const dx = footX - x;
        const dy = footY - y;

        if (Math.sqrt(dx * dx + dy * dy) < length)
            onGround = true;
    };

    this.draw = context => {
        const dx = footX - x;
        const dy = footY - y;
        const dist = Math.sqrt(dx * dx + dy * dy);
        const elbowAngle = Math.acos(dist / length) * Math.sign(direction);
        const footDirection = Math.atan2(footY - y, footX - x);
        const elbowX = x + Math.cos(footDirection + elbowAngle) * length * 0.5;
        const elbowY = y + Math.sin(footDirection + elbowAngle) * length * 0.5;

        context.fillStyle = "gray";
        context.strokeStyle = "black";

        context.beginPath();
        context.arc(footX, footY, 6, 0, Math.PI * 2);
        context.fill();
        context.stroke();

        context.beginPath();
        context.moveTo(x, y);
        context.lineTo(elbowX, elbowY);
        context.lineTo(footX, footY);
        context.stroke();
    };

    this.update = (newX, newY, bugDirection, speed, timeStep) => {
        x = newX;
        y = newY;

        const dx = footX - x;
        const dy = footY - y;
        const dist = Math.sqrt(dx * dx + dy * dy);

        if (onGround) {
            if (dist > length) {
                onGround = false;

                if (counterpart)
                    counterpart.stepDown();
            }
        }
        else {
            const xAim = x + Math.cos(bugDirection + direction) * length;
            const yAim = y + Math.sin(bugDirection + direction) * length;
            const dxAim = xAim - footX;
            const dyAim = yAim - footY;
            const lengthAim = Math.sqrt(dxAim * dxAim + dyAim * dyAim);

            if (lengthAim < length * (1 - Leg.GROUND_THRESHOLD))
                onGround = true;
            else {
                footX += (dxAim / lengthAim) * speed * timeStep;
                footY += (dyAim / lengthAim) * speed * timeStep;
            }
        }
    };

    applyInitialProgress();
};

Leg.GROUND_THRESHOLD = 0.9;
