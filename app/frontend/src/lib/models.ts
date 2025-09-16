//Models for go types.
export class Device {
  path: number[];
  name: string;
  vendorId: string;
  productId: string;
  speed: string;
  bus: number;
  state: string;

  static createFrom(source: any = {}) {
    return new Device(source);
  }

  constructor(source: any = {}) {
    if ("string" === typeof source) source = JSON.parse(source);
    this.path = source["path"];
    this.name = source["name"];
    this.vendorId = source["vendorId"];
    this.productId = source["productId"];
    this.speed = source["speed"];
    this.bus = source["bus"];
    this.state = source["state"];
  }
}

export class TreeNode {
  device: Device;
  children: TreeNode[];

  static createFrom(source: any = {}) {
    return new TreeNode(source);
  }

  constructor(source: any = {}) {
    if ("string" === typeof source) source = JSON.parse(source);
    this.device = source["device"];
    this.children = this.convertValues(source["children"], TreeNode);
  }

  convertValues(a: any, classs: any, asMap: boolean = false): any {
    if (!a) {
      return a;
    }
    if (a.slice && a.map) {
      return (a as any[]).map((elem) => this.convertValues(elem, classs));
    } else if ("object" === typeof a) {
      if (asMap) {
        for (const key of Object.keys(a)) {
          a[key] = new classs(a[key]);
        }
        return a;
      }
      return new classs(a);
    }
    return a;
  }
}

export class Log {
  // Go type: time
  Time: Date;
  Text: string;
  State: string;
  Speed: string;

  static createFrom(source: any = {}) {
    return new Log(source);
  }

  constructor(source: any = {}) {
    if ("string" === typeof source) source = JSON.parse(source);
    this.Time = this.convertValues(source["Time"], null);
    this.Text = source["Text"];
    this.State = source["State"];
    this.Speed = source["Speed"];
  }

  convertValues(a: any, classs: any, asMap: boolean = false): any {
    if (!a) {
      return a;
    }
    if (a.slice && a.map) {
      return (a as any[]).map((elem) => this.convertValues(elem, classs));
    } else if ("object" === typeof a) {
      if (asMap) {
        for (const key of Object.keys(a)) {
          a[key] = new classs(a[key]);
        }
        return a;
      }
      return new classs(a);
    }
    return a;
  }
}
