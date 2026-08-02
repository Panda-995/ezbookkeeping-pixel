import { readFileSync } from "node:fs";

import { parse } from "@vue/compiler-sfc";
import { describe, expect, it } from "vitest";

interface TemplateAstNode {
    readonly type: number;
    readonly tag?: string;
    readonly children?: readonly unknown[];
}

function isTemplateAstNode(node: unknown): node is TemplateAstNode {
    return typeof node === "object" && node !== null && "type" in node;
}

function findElement(node: TemplateAstNode, tag: string): TemplateAstNode | null {
    for (const child of node.children || []) {
        if (!isTemplateAstNode(child) || child.type !== 1) {
            continue;
        }

        if (child.tag === tag) {
            return child;
        }

        const matchingDescendant = findElement(child, tag);

        if (matchingDescendant) {
            return matchingDescendant;
        }
    }

    return null;
}

describe("desktop main layout route transition", () => {
    it("wraps routed components in an element that Vue can transition", () => {
        const source = readFileSync(
            new URL("./MainLayout.vue", import.meta.url),
            "utf8",
        );
        const template = parse(source).descriptor.template;

        if (!template?.ast) {
            throw new Error("MainLayout.vue must contain a template");
        }

        const transition = findElement(template.ast, "transition");

        expect(transition).not.toBeNull();

        const transitionRootElements = (transition!.children || []).filter(
            (child): child is TemplateAstNode =>
                isTemplateAstNode(child) && child.type === 1,
        );

        expect(transitionRootElements).toHaveLength(1);
        expect(transitionRootElements[0]?.tag).toBe("div");
        expect(findElement(transitionRootElements[0]!, "component")).not.toBeNull();
    });
});
