export default class Player {
    animationFrame;

    constructor(id, name, position, angle, ctx) {
        this.id = id;
        this.name = name;
        this.position = position;
        this.width = 50;
        this.height = 50;
        this.angle = angle;
        
        this.ctx = ctx;
    }

    draw() {
        // this.ctx.save();

        this.ctx.beginPath();
        this.ctx.fillStyle = "red";
        console.log("DRAW", this.position.x, this.position.y, this.width, this.height, this.angle);
        this.ctx.translate(this.position.x + this.width/2, this.position.y + this.height/2 );
        this.ctx.rotate(this.angle);
        this.ctx.fillRect(-this.width/2, -this.height/2, this.width, this.height);
        // this.ctx.closePath();
        

        // Reset current transformation matrix to the identity matrix
        this.ctx.setTransform(1, 0, 0, 1, 0, 0);

        this.ctx.fillStyle = "#ffffff";
        this.ctx.font = "18px Arial";
        this.ctx.textAlign = "center";
        this.ctx.textBaseline = "middle";
        this.ctx.fillText(this.name, this.position.x + this.width / 2, this.position.y + this.height / 2);

        this.ctx.closePath();


        // this.ctx.restore();
    }
}