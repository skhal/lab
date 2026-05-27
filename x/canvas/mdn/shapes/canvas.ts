// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// 2D coordinates use inverted y-scale:
// +----> X
// |
// |
// v
//
// Y
//
// coordinates use pixel units.

function main() {
  const canvas = <HTMLCanvasElement> document.getElementById("canvas");
  if (canvas.getContext == null) {
    // canvas is not supported
    return;
  }
  const ctx = canvas.getContext("2d");
  if (ctx == null) {
    return;
  }

  const rect = new RectDemo(ctx);
  rect.Stroke();
  rect.Fill();

  const line = new PathDemo(ctx);
  line.Stroke();
  line.Fill();

  const arc = new ArcDemo(ctx);
  arc.Stroke();
  arc.Fill();
}

class canvas {
  private readonly ctx: CanvasRenderingContext2D;
  constructor(ctx: CanvasRenderingContext2D) {
    this.ctx = ctx;
  }
  protected drawAt(x: number, y: number, f: any) {
    this.ctx.save();
    this.ctx.translate(x, y);
    f(this.ctx);
    this.ctx.restore();
  }
}

class RectDemo extends canvas {
  private path: Path2D = function (): Path2D {
    const p = new Path2D();
    p.rect(10, 10, 10, 10);
    return p;
  }();
  private rect: Path2D = function (): any {
    const p = new Path2D();
    p.rect(10, 10, 50, 50);
    return p;
  }();

  Stroke() {
    const that = this;
    // the line is blurry because it falls at the edge between actual pixels
    // https://developer.mozilla.org/en-US/docs/Web/API/Canvas_API/Tutorial/Drawing_shapes
    this.drawAt(0, 0, function (ctx: CanvasRenderingContext2D) {
      ctx.stroke(that.path);
    });

    // This time shift by .5 to position coordinates in the middle of the pixel.
    this.drawAt(20.5, 0.5, function (ctx: CanvasRenderingContext2D) {
      ctx.stroke(that.path);
    });
  }

  Fill() {
    const that = this;
    this.drawAt(0, 20, function (ctx: CanvasRenderingContext2D) {
      ctx.fillStyle = "rgb(200 0 0)";
      ctx.fill(that.rect);
    });

    this.drawAt(20, 40, function (ctx: CanvasRenderingContext2D) {
      ctx.fillStyle = "rgb(0 0 200 / 50%)";
      ctx.fill(that.rect);
    });
  }
}

class PathDemo extends canvas {
  private path = function (): Path2D {
    const p = new Path2D();
    p.moveTo(10, 10);
    p.lineTo(50, 10);
    p.lineTo(30, 30);
    p.closePath();
    return p;
  }();

  Stroke() {
    const that = this;
    this.drawAt(80, 0, function (ctx: CanvasRenderingContext2D) {
      ctx.strokeStyle = "rgb(0 200 0)";
      ctx.stroke(that.path);
    });
  }

  Fill() {
    const that = this;
    this.drawAt(80, 30, function (ctx: CanvasRenderingContext2D) {
      ctx.fillStyle = "rgb(0 0 200)";
      ctx.fill(that.path);
    });
  }
}

class ArcDemo extends canvas {
  private strokeStyle: string = "rgb(200 0 0)";
  private fillStyle: string = "rgb(0 0 200 / 50%)";
  private arcParams: any = {
    x: 20,
    y: 20,
    r: 10,
    start: 0,
    end: Math.PI / 2,
    ccw: false,
  };
  private arcClockwise = function (that: ArcDemo): Path2D {
    const p = new Path2D();
    const { x, y, r, start, end, ccw } = that.arcParams;
    p.arc(x, y, r, start, end, ccw);
    return p;
  }(this);
  private arcCounterClockwise = function (that: ArcDemo): Path2D {
    const p = new Path2D();
    const { x, y, r, start, end, ccw } = {
      ...that.arcParams,
      ...{ ccw: true },
    };
    p.arc(x, y, r, start, end, ccw);
    return p;
  }(this);

  Stroke() {
    const that = this;
    this.drawAt(80, 60, function (ctx: CanvasRenderingContext2D) {
      ctx.strokeStyle = that.strokeStyle;
      ctx.stroke(that.arcCounterClockwise);
    });
    this.drawAt(110, 60, function (ctx: CanvasRenderingContext2D) {
      ctx.strokeStyle = that.strokeStyle;
      ctx.stroke(that.arcClockwise);
    });
  }

  Fill() {
    const that = this;
    this.drawAt(80, 90, function (ctx: CanvasRenderingContext2D) {
      ctx.fillStyle = that.fillStyle;
      ctx.fill(that.arcCounterClockwise);
    });
    this.drawAt(110, 90, function (ctx: CanvasRenderingContext2D) {
      ctx.fillStyle = that.fillStyle;
      ctx.fill(that.arcClockwise);
    });
  }
}
