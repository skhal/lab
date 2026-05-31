// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

type appConfig = {
  canvas: string;
  buttons: {
    start: {
      id: string;
    };
    stop: {
      id: string;
    };
  };
};

export class App {
  private cfg: {
    canvas: HTMLCanvasElement;
    buttons: {
      start: HTMLButtonElement;
      stop: HTMLButtonElement;
    };
  };

  constructor(cfg: appConfig) {
    this.cfg = {
      canvas: <HTMLCanvasElement> document.getElementById(cfg.canvas),
      buttons: {
        start: <HTMLButtonElement> document.getElementById(
          cfg.buttons.start.id,
        ),
        stop: <HTMLButtonElement> document.getElementById(cfg.buttons.stop.id),
      },
    };
    const startOnclick = this.cfg.buttons.start.onclick;
    this.cfg.buttons.start.onclick = (ev) => {
      this.cfg.buttons.start.disabled = true;
      this.cfg.buttons.stop.disabled = false;
      startOnclick?.call(this.cfg.buttons.start, ev);
    };
    const stopOnclick = this.cfg.buttons.stop.onclick;
    this.cfg.buttons.stop.onclick = (ev) => {
      this.cfg.buttons.start.disabled = false;
      this.cfg.buttons.stop.disabled = true;
      stopOnclick?.call(this.cfg.buttons.stop, ev);
    };
  }

  public Run() {
    this.cfg.buttons.start.disabled = false;
  }
}
