package compiler_objects

type UnverifiedElementQueue struct {
	elements []UnverifiedElement
}

func (queue *UnverifiedElementQueue) enqueue(element UnverifiedElement) {
	queue.elements = append(queue.elements, element)
}

func (queue *UnverifiedElementQueue) dequeue() UnverifiedElement {
	if len(queue.elements) < 1 {
		return UnverifiedElement{}
	}
	element := queue.elements[0]
	queue.elements = queue.elements[1:]
	return element
}
