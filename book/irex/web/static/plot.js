// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Init initializes animation in SVG plot with a given id.
export function Init(id, data) {
    const svg = document.getElementById(id);
    if (svg == null) {
        throw new Error(`can't get SVG element #${id}`);
    }
    const p = new plotter(svg, data);
    p.Run();
}
class plotter {
    svg;
    lineGroup;
    linePath;
    lineRect;
    pointerGroup;
    pointerGuide;
    pointerCursor;
    hintVal;
    hintDate;
    coords;
    svgPoint;
    lastX = 0;
    data;
    constructor(svg, data) {
        this.svg = svg;
        this.data = data;
        this.lineGroup = svg.getElementById("plot-line-group");
        if (this.lineGroup == null) {
            throw new Error(`svg ${this.svg.getAttribute("id")}: missing #plot-line-group`);
        }
        this.linePath = svg.getElementById("plot-line-path");
        if (this.linePath == null) {
            throw new Error(`svg ${this.svg.getAttribute("id")}: missing #plot-line-path`);
        }
        this.lineRect = svg.getElementById("plot-line-rect");
        if (this.lineRect == null) {
            throw new Error(`svg ${this.svg.getAttribute("id")}: missing #plot-line-rect`);
        }
        this.pointerGroup = svg.getElementById("plot-line-pointer");
        if (this.pointerGroup == null) {
            throw new Error(`svg ${this.svg.getAttribute("id")}: missing #plot-line-pointer`);
        }
        this.pointerGuide = svg.getElementById("plot-line-pointer-guide");
        if (this.pointerGuide == null) {
            throw new Error(`svg ${this.svg.getAttribute("id")}: missing #plot-line-guide`);
        }
        this.pointerCursor = svg.getElementById("plot-line-pointer-cursor");
        if (this.pointerCursor == null) {
            throw new Error(`svg ${this.svg.getAttribute("id")}: missing #plot-line-pointer-cursor`);
        }
        this.hintVal = document.getElementById("plot-hint-val");
        if (this.hintVal == null) {
            throw new Error(`missing #plot-hint-val`);
        }
        this.hintDate = document.getElementById("plot-hint-date");
        if (this.hintDate == null) {
            throw new Error(`missing #plot-hint-date`);
        }
        this.coords = coordinatesFrom(this.linePath);
        this.svgPoint = this.svg.createSVGPoint();
    }
    Run() {
        this.lineRect.addEventListener("mousemove", (e) => this.lineRectMouseMove(e));
    }
    lineRectMouseMove(e) {
        this.svgPoint.x = e.clientX;
        this.svgPoint.y = e.clientY;
        const loc = this.svgPoint.matrixTransform(this.lineGroup.getScreenCTM()?.inverse());
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
function coordinatesFrom(el) {
    const m = new Map();
    for (let d = 0; d < el.getTotalLength(); d++) {
        const p = el.getPointAtLength(d);
        const x = Math.floor(p.x);
        m.set(x, p.y);
    }
    return m;
}
