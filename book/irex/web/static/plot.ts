// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Init initializes animation in SVG plot with a given id.
export function Init(id: string, data: any) {
  const svg = <SVGSVGElement> <any> document.getElementById(id);
  if (svg == null) {
    throw new Error(`can't get SVG element #${id}`);
  }
  const p = new plotter(svg, data);
  p.Run();
}

class plotter {
  private svg: SVGSVGElement;
  private lineGroup: SVGGElement;
  private linePath: SVGPathElement;
  private lineRect: SVGRectElement;
  private pointerGroup: SVGGElement;
  private pointerGuide: SVGLineElement;
  private pointerCursor: SVGCircleElement;
  private hintVal: HTMLSpanElement;
  private hintDate: HTMLSpanElement;
  private coords: Map<number, number>;
  private svgPoint: SVGPoint;
  private lastX: number = 0;
  private data: any;

  constructor(svg: SVGSVGElement, data: any) {
    this.svg = svg;
    this.data = data;
    this.lineGroup = <SVGGElement> svg.getElementById("plot-line-group");
    if (this.lineGroup == null) {
      throw new Error(
        `svg ${this.svg.getAttribute("id")}: missing #plot-line-group`,
      );
    }
    this.linePath = <SVGPathElement> svg.getElementById("plot-line-path");
    if (this.linePath == null) {
      throw new Error(
        `svg ${this.svg.getAttribute("id")}: missing #plot-line-path`,
      );
    }
    this.lineRect = <SVGRectElement> svg.getElementById("plot-line-rect");
    if (this.lineRect == null) {
      throw new Error(
        `svg ${this.svg.getAttribute("id")}: missing #plot-line-rect`,
      );
    }
    this.pointerGroup = <SVGGElement> svg.getElementById("plot-line-pointer");
    if (this.pointerGroup == null) {
      throw new Error(
        `svg ${this.svg.getAttribute("id")}: missing #plot-line-pointer`,
      );
    }
    this.pointerGuide = <SVGLineElement> svg.getElementById(
      "plot-line-pointer-guide",
    );
    if (this.pointerGuide == null) {
      throw new Error(
        `svg ${this.svg.getAttribute("id")}: missing #plot-line-guide`,
      );
    }
    this.pointerCursor = <SVGCircleElement> svg.getElementById(
      "plot-line-pointer-cursor",
    );
    if (this.pointerCursor == null) {
      throw new Error(
        `svg ${this.svg.getAttribute("id")}: missing #plot-line-pointer-cursor`,
      );
    }
    this.hintVal = <HTMLSpanElement> document.getElementById(
      "plot-hint-val",
    );
    if (this.hintVal == null) {
      throw new Error(
        `missing #plot-hint-val`,
      );
    }
    this.hintDate = <HTMLSpanElement> document.getElementById(
      "plot-hint-date",
    );
    if (this.hintDate == null) {
      throw new Error(
        `missing #plot-hint-date`,
      );
    }
    this.coords = coordinatesFrom(this.linePath);
    this.svgPoint = this.svg.createSVGPoint();
  }

  public Run() {
    this.lineRect.addEventListener(
      "mousemove",
      (e) => this.lineRectMouseMove(e),
    );
  }

  private lineRectMouseMove(e: MouseEvent) {
    this.svgPoint.x = e.clientX;
    this.svgPoint.y = e.clientY;
    const loc = this.svgPoint.matrixTransform(
      this.lineGroup.getScreenCTM()?.inverse(),
    );
    if (loc == null) {
      throw new Error("failed to get translation matrix");
    }
    const x = Math.floor(loc.x);
    if (this.lastX == x) {
      return;
    }
    this.lastX = x;
    const y = this.coords.get(x);
    if (y == undefined) {
      return;
    }
    this.pointerGroup.setAttribute("transform", `translate(${x} 0)`);
    this.pointerCursor.setAttribute("cy", y.toString());
    this.pointerGroup.style.visibility = "visible";
    const v = this.data.get(x);
    if (v == undefined) {
      return;
    }
    const val = parseInt(v.c);
    if (isNaN(val)) {
      return;
    }
    const unix_ts = parseInt(v.d);
    if (isNaN(unix_ts)) {
      return;
    }
    const quotient = Math.floor(val / 100);
    const remainder = val % 100;
    const date = new Date(unix_ts * 1000); // must be milliseconds
    this.hintDate.textContent = date.toLocaleDateString(undefined, {
      year: "numeric",
      month: "short",
    });
    this.hintVal.textContent = `${quotient}.${remainder}`;
  }
}

function coordinatesFrom(el: SVGPathElement): Map<number, number> {
  const m = new Map();
  for (let d = 0; d < el.getTotalLength(); d++) {
    const p = el.getPointAtLength(d);
    const x = Math.floor(p.x);
    m.set(x, p.y);
  }
  return m;
}
