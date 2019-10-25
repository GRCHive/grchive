export function convertClientDeltaToSvg(svg : SVGSVGElement, dx : number, dy : number) : SVGPoint {
    let pt : SVGPoint = svg.createSVGPoint()
    pt.x = dx
    pt.y = dy

    let matrix : DOMMatrix = svg.getScreenCTM()!
    matrix.e = 0
    matrix.f = 0

    return pt.matrixTransform(matrix.inverse())
}

export function convertClientPointToSvg(svg : SVGSVGElement, x : number, y : number) : SVGPoint {
    let pt : SVGPoint = svg.createSVGPoint()
    pt.x = x
    pt.y = y
    return pt.matrixTransform(svg.getScreenCTM()!.inverse())
}
