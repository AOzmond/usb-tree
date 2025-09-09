export namespace lib {
  export class Device {
    Path: number[];
    Name: string;
    VendorID: string;
    ProductID: string;
    Speed: string;
    Bus: number;
    State: string;

    static createFrom(source: any = {}) {
      return new Device(source);
    }

    constructor(source: any = {}) {
      if ("string" === typeof source) source = JSON.parse(source);
      this.Path = source["Path"];
      this.Name = source["Name"];
      this.VendorID = source["VendorID"];
      this.ProductID = source["ProductID"];
      this.Speed = source["Speed"];
      this.Bus = source["Bus"];
      this.State = source["State"];
    }
  }
}
