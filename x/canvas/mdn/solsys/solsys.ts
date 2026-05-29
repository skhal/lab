// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

export function Main(canvasID: string) {
  const canvas = <HTMLCanvasElement> document.getElementById(canvasID);
  if (canvas == null) {
    console.error(`${canvasID} does not exist or is not a canvas`);
    return;
  }
  const ctx = canvas.getContext("2d");
  if (ctx == null) {
    console.error(`can't get ${canvasID} 2D context`);
    return;
  }
  const d = new drawer(ctx, { width: canvas.width, height: canvas.height });
  d.Run();
}

type dimensions = {
  width: number;
  height: number;
};

class drawer {
  private ctx: CanvasRenderingContext2D;
  private dim: dimensions;
  private sun: Sun;
  private earth: Earth;
  private moon: Moon;

  constructor(ctx: CanvasRenderingContext2D, dim: dimensions) {
    this.ctx = ctx;
    this.dim = dim;
    this.sun = new Sun(ctx);
    this.earth = new Earth(ctx, dim.width / 4);
    this.moon = new Moon(ctx, this.dim.width / 18);
  }

  public Run() {
    Promise.all([
      this.sun.Ready(),
      this.earth.Ready(),
      this.moon.Ready(),
    ]).then(() => this.draw());
  }

  private draw() {
    this.render();
    window.requestAnimationFrame(() => this.draw());
  }

  private render() {
    stateLock(this.ctx, () => this.reset());
    stateLock(this.ctx, () => this.renderSun());
  }

  private reset() {
    this.ctx.fillStyle = "black";
    this.ctx.fillRect(0, 0, this.dim.width, this.dim.height);
  }

  private renderSun() {
    this.ctx.translate(this.dim.width / 2, this.dim.height / 2);
    this.sun.Draw();
    this.earth.Draw();
    this.moon.Draw();
  }
}

// image wraps HTMLImageElement to start loading the image in the constructor
// and let clients use [Ready] to wait for the image to be ready.
class image {
  elem: HTMLImageElement = new Image();
  dim: dimensions;

  private load: Promise<HTMLImageElement>;

  constructor(url: string, dim: dimensions) {
    this.load = new Promise((resolve) => {
      const img = new Image();
      img.addEventListener("load", () => resolve(img));
      img.src = url;
    });
    this.dim = dim;
  }

  public async Ready() {
    this.elem = await this.load;
  }
}

const animationFPS = 60;

function angleDelta(spinFPS: number): number {
  const fps = spinFPS * animationFPS;
  return 2 * Math.PI / fps;
}

type astronomicalObjectConfig = {
  orbitRadius: number;
  orbitSeconds: number;
  spinSeconds: number;
};

class astronomicalObject {
  private frame: number = 0;
  private cfg: {
    orbit: { r: number; dphi: number };
    spin: { dphi: number };
  };

  protected ctx: CanvasRenderingContext2D;

  constructor(
    ctx: CanvasRenderingContext2D,
    cfg: astronomicalObjectConfig,
  ) {
    this.ctx = ctx;
    this.cfg = {
      orbit: { r: cfg.orbitRadius, dphi: angleDelta(cfg.orbitSeconds) },
      spin: { dphi: angleDelta(cfg.spinSeconds) },
    };
  }

  public Draw() {
    this.frame++;
    if (this.cfg.orbit.r != 0) {
      stateLock(this.ctx, () => this.drawOrbit());
    }
    this.position();
    this.spin();
  }

  private drawOrbit() {
    this.ctx.strokeStyle = "rgb(146 146 146/ 50%)";
    this.ctx.lineWidth = 2;
    this.ctx.setLineDash([4, 8]);
    this.ctx.beginPath();
    this.ctx.moveTo(this.cfg.orbit.r, 0);
    this.ctx.arc(0, 0, this.cfg.orbit.r, 0, 2 * Math.PI, true);
    this.ctx.stroke();
  }

  private position() {
    const angle = this.frame * this.cfg.orbit.dphi;
    this.ctx.rotate(angle);
    this.ctx.translate(this.cfg.orbit.r, 0);
  }

  private spin() {
    const angle = this.frame * this.cfg.spin.dphi;
    this.ctx.rotate(angle);
  }
}

class Sun extends astronomicalObject {
  private static config = {
    radius: 50,
    orbitRadius: 0,
    orbitSeconds: 0,
    spinSeconds: 120,
  };

  private img: image;

  constructor(ctx: CanvasRenderingContext2D) {
    super(ctx, Sun.config);
    const dim = {
      width: 2 * Sun.config.radius,
      height: 2 * Sun.config.radius,
    };
    this.img = new image("img/sun.jpg", dim);
  }

  public async Ready() {
    return this.img.Ready();
  }

  public Draw() {
    super.Draw();
    stateLock(this.ctx, () => this.draw());
  }

  private draw() {
    this.ctx.beginPath();
    this.ctx.arc(0, 0, this.img.dim.width / 2, 0, 2 * Math.PI, true);
    this.ctx.clip();

    this.ctx.translate(-this.img.dim.width / 2, -this.img.dim.height / 2);
    this.ctx.drawImage(
      this.img.elem,
      0,
      0,
      this.img.dim.width,
      this.img.dim.height,
    );
  }
}

class Earth extends astronomicalObject {
  private static config = {
    radius: 25,
    orbitSeconds: 30,
    spinSeconds: 5,
  };

  private img: image;

  constructor(ctx: CanvasRenderingContext2D, orbit: number) {
    const cfg = { ...Earth.config, ...{ orbitRadius: orbit } };
    super(ctx, cfg);
    const dim = {
      width: 2 * Earth.config.radius,
      height: 2 * Earth.config.radius,
    };
    this.img = new image("img/earth.jpg", dim);
  }

  public async Ready() {
    return this.img.Ready();
  }

  public Draw() {
    super.Draw();
    stateLock(this.ctx, () => this.draw());
  }

  private draw() {
    this.ctx.beginPath();
    this.ctx.arc(0, 0, this.img.dim.width * .35, 0, 2 * Math.PI, true);
    this.ctx.clip();

    this.ctx.translate(-this.img.dim.width / 2, -this.img.dim.height / 2);
    this.ctx.drawImage(
      this.img.elem,
      0,
      0,
      this.img.dim.width,
      this.img.dim.height,
    );
  }
}

class Moon extends astronomicalObject {
  private static config = {
    radius: 5,
    orbitSeconds: 10,
    spinSeconds: 1,
  };

  private img: image;

  constructor(ctx: CanvasRenderingContext2D, orbit: number) {
    const cfg = { ...Moon.config, ...{ orbitRadius: orbit } };
    super(ctx, cfg);
    const dim = {
      width: 2 * Moon.config.radius,
      height: 2 * Moon.config.radius,
    };
    this.img = new image("img/moon.jpg", dim);
  }

  public async Ready() {
    return this.img.Ready();
  }

  public Draw() {
    super.Draw();
    stateLock(this.ctx, () => this.draw());
  }

  private draw() {
    this.ctx.beginPath();
    this.ctx.arc(0, 0, this.img.dim.width / 2, 0, 2 * Math.PI, true);
    this.ctx.clip();

    this.ctx.translate(-this.img.dim.width / 2, -this.img.dim.height / 2);
    this.ctx.drawImage(
      this.img.elem,
      0,
      0,
      this.img.dim.width,
      this.img.dim.height,
    );
  }
}

function stateLock(ctx: CanvasRenderingContext2D, f: () => void) {
  ctx.save();
  f();
  ctx.restore();
}
