const pointsData = [
	{x: 0, y: -0.2, parent: -1}, // head
	{x: -0.3, y: -0.1, parent: -1}, // left shoulder (from screen view)
	{x: -0.2, y: -0.4, parent: 1},
	{x: 0.2, y: -0.1, parent: -1}, // right shoulder
	{x: 0.4, y: -0.4, parent: 3},
	{x: 0.3, y: 0.1, parent: -1}, // right thigh
	{x: 0.4, y: 0.2, parent: 5},
	{x: -0.1, y: 0.2, parent: -1}, // left thigh
	{x: 0, y: 0.4, parent: 7},
]
const initialPoints = [];
const points = [];
const stems = [3, 5, 7, 1, 0]
const parent = pointsData.map(e => e.parent);

let captured = -1;

function setup() {
  let w = windowWidth;
  // if (windowWidth > 640) {
  //   w = 480;
  // }
	const c = document.querySelector("#p5container")
	createCanvas(c?.offsetWidth, c?.offsetHeight, c?.querySelector("canvas"));
	
	initialPoints.push(...pointsData.map(pt => createVector(pt.x, pt.y)))
	points.push(...pointsData.map(pt => createVector(pt.x, pt.y)))
}

function windowResized() {
  let w = windowWidth;
  // if (windowWidth > 640) {
  //   w = 480;
  // }
	const c = document.querySelector("#p5container")
	resizeCanvas(c?.offsetWidth, c?.offsetHeight);
}

function draw() {
	clear();
	const T = millis() * 0.001;
	curveTightness(0);
	const R = Math.min(width, height);
	translate(width / 2, height / 2);
	
	for (const i in points) {
		const nv = p5.Vector.fromAngle(((noise(T, i * 0.1)-0.5) * Math.PI * 2));
		nv.setMag((noise(i * 0.1, T) - 0.5) * 0.001)
		initialPoints[i].add(nv);
		
		const m = 0.45;
		if (initialPoints[i].x > m) initialPoints[i].x = m;
		if (initialPoints[i].x < -m) initialPoints[i].x = -m;
		if (initialPoints[i].y > m) initialPoints[i].y = m;
		if (initialPoints[i].y < -m) initialPoints[i].y = -m;
		
		const v = p5.Vector.sub(points[i], initialPoints[i]);
		v.limit(0.001);
		points[i].sub(v);
	}
	
	const centroid = stems.reduce((acc, cur) => p5.Vector.add(acc, points[cur]), createVector(0, 0))
	centroid.div(stems.length);
	
	push(); {
		const gradient = drawingContext
		  .createRadialGradient(centroid.x * R, centroid.y * R, R * 0.1,
														points[0].x * R, points[0].y * R, R);

		// Add three color stops
		gradient.addColorStop(0, `hsl(${(T*10)%360}deg 100% 50%)`);
		// gradient.addColorStop(0.5, "blue");
		gradient.addColorStop(1, `hsl(${(T*10+90)%360}deg 100% 50%)`);

		// Set the fill style and draw a rectangle
		drawingContext.fillStyle = gradient;

		beginShape();
		let first;
		let second;
		let third;
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
				first = createVector((pt.x + vv.x * 0.1) * R,
														 (pt.y + vv.y * 0.1) * R);
			}
			vv.rotate(Math.PI/2);
			if (childP !== undefined) {
				curveVertex((childP.x + vc.x * 0.1) * R,
										(childP.y + vc.y * 0.1) * R);			
				if (second === undefined) {
					second = createVector((childP.x + vc.x * 0.1) * R,
										(childP.y + vc.y * 0.1) * R);
				}
				vc.rotate(Math.PI/2)
				curveVertex((childP.x + vc.x * 0.1) * R,
										(childP.y + vc.y * 0.1) * R);			
				if (third === undefined) {
					third = createVector((childP.x + vc.x * 0.1) * R,
										(childP.y + vc.y * 0.1) * R);
				}
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
		curveVertex(first.x, first.y);
		curveVertex(second.x, second.y);
		curveVertex(third.x, third.y);
		endShape();
	} pop();

	for (const pt of points) {
		noStroke();
		const x = pt.x * R;
		const y = pt.y * R;
  	circle(x, y, R*0.03);
	}
	push(); {
		noStroke();
		const x = centroid.x * R;
		const y = centroid.y * R;
		translate(x, y);
		const r = R * 0.07 * map(Math.sin(T * Math.PI * 2 * 0.1), -1, 1, 0.5, 1);
		fill(0, 200)
  	circle(R*0.001, R*0.001, r*1.002);
		fill("red")
  	circle(0, 0, r);
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
