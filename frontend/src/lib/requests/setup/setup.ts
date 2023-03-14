// @generated by protobuf-ts 2.8.2
// @generated from protobuf file "setup/setup.proto" (package "setup", syntax proto3)
// tslint:disable
import type { BinaryWriteOptions } from "@protobuf-ts/runtime";
import type { IBinaryWriter } from "@protobuf-ts/runtime";
import { WireType } from "@protobuf-ts/runtime";
import type { BinaryReadOptions } from "@protobuf-ts/runtime";
import type { IBinaryReader } from "@protobuf-ts/runtime";
import { UnknownFieldHandler } from "@protobuf-ts/runtime";
import type { PartialMessage } from "@protobuf-ts/runtime";
import { reflectionMergePartial } from "@protobuf-ts/runtime";
import { MESSAGE_TYPE } from "@protobuf-ts/runtime";
import { MessageType } from "@protobuf-ts/runtime";
import { Family } from "../core/family";
import { User } from "../core/user";
/**
 * @generated from protobuf message setup.Setup
 */
export interface Setup {
    /**
     * @generated from protobuf field: core.User user = 1;
     */
    user?: User;
    /**
     * @generated from protobuf field: core.Family family = 2;
     */
    family?: Family;
}
// @generated message type with reflection information, may provide speed optimized methods
class Setup$Type extends MessageType<Setup> {
    constructor() {
        super("setup.Setup", [
            { no: 1, name: "user", kind: "message", T: () => User },
            { no: 2, name: "family", kind: "message", T: () => Family }
        ]);
    }
    create(value?: PartialMessage<Setup>): Setup {
        const message = {};
        globalThis.Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial<Setup>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Setup): Setup {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* core.User user */ 1:
                    message.user = User.internalBinaryRead(reader, reader.uint32(), options, message.user);
                    break;
                case /* core.Family family */ 2:
                    message.family = Family.internalBinaryRead(reader, reader.uint32(), options, message.family);
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: Setup, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* core.User user = 1; */
        if (message.user)
            User.internalBinaryWrite(message.user, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* core.Family family = 2; */
        if (message.family)
            Family.internalBinaryWrite(message.family, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message setup.Setup
 */
export const Setup = new Setup$Type();