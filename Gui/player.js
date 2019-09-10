class Sector {

    getX() {
        return this.x;
    }
    getY() {
        return this.y;
    }

    setX(x) {
        this.x = x;
    }
    setY(y) {
        this.y = y;
    }

    diffToAngle(x,y) {
        let angle = Math.atan2(y, x) - Math.atan2(0, 1);
        if (angle < 0) {
            angle += 2 * Math.PI;
        }
        return angle;
    }

    update(x, y) {
        this.diffX = x - this.x;
        this.diffY = y - this.y;

        this.x = x;
        this.y = y;

        this.dirAngle = this.diffToAngle(this.diffX, this.diffY);

        let speed = 60.5;
        let timeStep = 0.017;

        for (let leg of this.legs) {
            leg.update(this.x, this.y, this.dirAngle, speed * 4, timeStep);
        }
    }

    draw(context) {
        for (let leg of this.legs) {
            leg.draw(context);
        }
    }

    constructor(x, y, createLegs, body) {
        this.x = x;
        this.y = y;
        this.legs = [];
        this.diffX = 0.0;
        this.diffY = 0.0;
        this.dirAngle = 0.0;

        if (createLegs) {
            let direction = 0.0;
            let speed = 10.5;
            const l = new Leg(x, y, direction, -body.getLegAngle(), body.getLegLength(), Math.random(), speed * 4);
            const r = new Leg(x, y, direction, body.getLegAngle(), body.getLegLength(), Math.random(), speed * 4);

            l.setCounterpart(r);
            r.setCounterpart(l);

            this.legs.push(l, r);
        }

    }
}

class Player {

    updateSize(size) {
        this.size = size;
    }

    updateHighscoreIndex(index) {
        this.highScoreIndex = index;
    }


    updatePositions(positions) {


        if (this.positions.length < positions.length) {
            this.positions = new Array(positions.length);
            for (let i = 0; i < positions.length; i++) {
                this.positions[i] = new Sector(positions[i][0], positions[i][1], positions.length == 1 || (i+2) % 2 == 0, this.body);
                //this.path.add(new paper.Point(positions[i][0], positions[i][1]));
            }
        }

        for (let i = 0; i < positions.length; i++) {
            this.positions[i].update(positions[i][0], positions[i][1]);
            //this.path.segments[i].point.x = positions[i][0];
            //this.path.segments[i].point.y = positions[i][1];
        }
    }

    drawBug(context) {
        for (let p of this.positions) {
            p.draw(context);
        }



    }

    updateBullets(bullets) {
        if (this.bullets.length != bullets.length) {
            this.bullets = new Array(bullets.length)
        }
        for (let i = 0; i < bullets.length; i++) {
            this.bullets[i] = bullets[i].pos;
        }
    }



    constructor(id, name, color, size, positions, bullets, hsIndex) {
        this.id = id;
        this.name = name;
        this.color = color;
        this.size = size;

        this.bullets = [];
        this.positions = [];
        this.highScoreIndex = hsIndex;

        this.body = new BodyShape();

        //this.path = new paper.Path({
        //    strokeColor: '#E4141B',
        //    strokeWidth: 20,
        //    strokeCap: 'round'
        //});

    }



}
