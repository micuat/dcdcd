const pointsData = [
	{x: 0, y: -0.2, parent: -1},
	{x: -0.3, y: -0.1, parent: -1},
	{x: -0.2, y: -0.4, parent: 1},
	{x: 0.2, y: -0.1, parent: -1},
	{x: 0.4, y: -0.4, parent: 3},
	{x: 0.3, y: 0.1, parent: -1},
	{x: 0.4, y: 0.2, parent: 5},
	{x: -0.1, y: 0.2, parent: -1},
	{x: 0, y: 0.4, parent: 7},
]
const points = [];
const stems = [1, 0, 3, 5, 7]
const parent = pointsData.map(e => e.parent);

let captured = -1;

function setup() {
  let w = windowWidth;
  if (windowWidth > 640) {
    w = 480;
  }
	createCanvas(w, windowHeight, document.querySelector("canvas"));
	
	points.push(...pointsData.map(pt => createVector(pt.x, pt.y)))
}

function windowResized() {
  let w = windowWidth;
  if (windowWidth > 640) {
    w = 480;
  }
  resizeCanvas(w, windowHeight);
}

function draw() {
	background(100);
	
	curveTightness(0);
	const R = Math.min(width, height);
	translate(width / 2, height / 2);
	
	const centroid = stems.reduce((acc, cur) => p5.Vector.add(acc, points[cur]), createVector(0, 0))
	centroid.div(stems.length);
	
	beginShape();
	let first;
	let last
	for (const idx of stems) {
		const pt = points[idx];
		const childIdx = parent.findIndex(e => e == idx);
		const childP = points[childIdx];
		const v = p5.Vector.sub(pt, centroid);
		// translate(centroid.x * R, centroid.y * R);
		// line(0, 0, v.x * R, v.y * R);
		const vv = v.copy();
		vv.normalize();
		let vc;
		if (childP !== undefined) {
			vc = p5.Vector.sub(childP, v);
			vc.normalize();
			const a0 = vv.heading();
			const a1 = vc.heading();
			vv.rotate((a1 - a0) / 2)
			vc.rotate(-Math.PI/2);
		}
		vv.rotate(-Math.PI/2);
		curveVertex((pt.x + vv.x * 0.1) * R,
								(pt.y + vv.y * 0.1) * R);
		if (first === undefined) {
		curveVertex((pt.x + vv.x * 0.1) * R,
								(pt.y + vv.y * 0.1) * R);
			first = createVector((pt.x + vv.x * 0.1) * R,
													 (pt.y + vv.y * 0.1) * R);
		}
		vv.rotate(Math.PI/2);
		if (childP !== undefined) {
			curveVertex((childP.x + vc.x * 0.1) * R,
									(childP.y + vc.y * 0.1) * R);			
			vc.rotate(Math.PI/2)
			curveVertex((childP.x + vc.x * 0.1) * R,
									(childP.y + vc.y * 0.1) * R);			
			vc.rotate(Math.PI/2)
			curveVertex((childP.x + vc.x * 0.1) * R,
									(childP.y + vc.y * 0.1) * R);			
		}
		else {
			curveVertex((pt.x + vv.x * 0.1) * R,
									(pt.y + vv.y * 0.1) * R);
		}
		vv.rotate(Math.PI/2);
		curveVertex((pt.x + vv.x * 0.1) * R,
								(pt.y + vv.y * 0.1) * R);
		last = createVector((pt.x + vv.x * 0.1) * R,
										 (pt.y + vv.y * 0.1) * R);

	}
	curveVertex(last.x, last.y);
	endShape();

	for (const pt of points) {
		const x = pt.x * R;
		const y = pt.y * R;
  	circle(x, y, R*0.01);
	}
	push(); {
		fill("red")
		const x = centroid.x * R;
		const y = centroid.y * R;
  	circle(x, y, 20);
	} pop();
}

function pressed() {
	const R = Math.min(width, height);
	const m = createVector((mouseX - width/2) / R, (mouseY - height/2) / R);
	for (const i in points) {
		if (points[i].dist(m) * R < 30) {
			captured = i;
			points[i].set(m.x, m.y);
			break;
		}
	}
}

function mousePressed() {
	pressed();
}

function touchStarted() {
	pressed();
}

function moved() {
	const R = Math.min(width, height);
	const m = createVector((mouseX - width/2) / R, (mouseY - height/2) / R);
	if (captured >= 0) {
		points[captured].set(m.x, m.y);
	}
}

function mouseDragged() {
	moved();
}

function touchMoved() {
	moved();
}

function released() {
	captured = -1;
}

function mouseReleased() {
	released();
}

function touchEnded() {
	released();
}