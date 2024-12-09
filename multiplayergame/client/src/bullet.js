export default class Bullet {

    constructor(position, angle, bulletType, radius, ctx) {
        this.position = position;
        this.angle = angle;
        this.bulletType = bulletType
        this.radius = radius;
        
        this.ctx = ctx;
    }

    draw() {
        this.ctx.beginPath();
        this.ctx.fillStyle = "black";
        this.ctx.arc(this.position.x, this.position.y, this.radius, 2 * Math.PI, false)
        this.ctx.fill();
        this.ctx.closePath();
    }
}