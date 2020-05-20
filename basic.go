package dominant_colour

import (
	"fmt"
	golang_queues "github.com/nadav-rahimi/golang-queue"
	gs "github.com/nadav-rahimi/golang-sets"
	"gonum.org/v1/gonum/mat"
	"math"
)

// Calculates the magnitude of a vector
func _3DVecMagnitude(v *mat.VecDense) float64 {
	var total float64 = 0
	total += math.Pow((v.AtVec(0)), 2)
	total += math.Pow((v.AtVec(1)), 2)
	total += math.Pow((v.AtVec(2)), 2)
	return math.Sqrt(total)
}

// Binary tree node used for the calculating the dominant colour
type tree_node struct {
	// Mean - Qn
	qn *mat.Dense
	// Covariance - Rn (with squiggly line over R)
	rn *mat.SymDense
	// Cardinality - Nn
	nn int
	// The set which contains the rgb pixels of the image
	pixels *gs.Set

	// Sum of all colours
	mn *mat.Dense
	// Rn (no squiggly line)
	rn_s *mat.Dense

	// Left node of the binary tree
	left *tree_node
	// Right node of the binary tree
	right *tree_node
}

// Returns a new tree node with a given pixel set
func newNode(pixels *gs.Set) *tree_node {
	return &tree_node{
		qn:     mat.NewDense(3, 1, nil),
		rn:     mat.NewSymDense(3, nil),
		nn:     pixels.Size(),
		pixels: pixels,
		mn:     mat.NewDense(3, 1, nil),
		rn_s:   mat.NewDense(3, 3, nil),
		left:   nil,
		right:  nil,
	}
}

// Wrapper for the mean and covariance functions to ensure one is called after the other
func (t *tree_node) calculate_mean_and_covariance() {
	t._calculate_mean()
	t._calculate_covariance()
}

// Calculates the mean of the tree node
func (t *tree_node) _calculate_mean() {
	// Calculates Mn
	t.mn.Zero()
	for _, v := range t.pixels.List() {
		matrix, ok := v.(mat.Matrix)
		chkFatal(ok, "retrieving pixels when calculating mean failed")
		t.mn.Add(t.mn, matrix)
	}

	// Calculates Qn by dividing Mn by the Cardinality of the Set
	t.qn = mat.NewDense(3, 1, nil)
	t.qn.Scale(float64(1)/float64(t.nn), t.mn)
}

// Calculates the covariance of the tree node
func (t *tree_node) _calculate_covariance() {
	// Calculating Rn (no squiggly line)
	var multi mat.Dense
	t.rn_s.Zero()

	for _, v := range t.pixels.List() {
		matrix, ok := v.(mat.Matrix)
		chkFatal(ok, "retrieving pixels when calculating covariance failed")

		multi.Mul(matrix, matrix.T())

		t.rn_s.Add(t.rn_s, &multi)
	}

	// Multiplying the mean by its transpose and scaling it down
	var mean_multi = mat.NewDense(3, 3, nil)
	mean_multi.Mul(t.qn, t.qn.T())

	// Minus the mean from Rn (non-squiggle) to get the covariance (R squiggly)
	rn := mat.NewDense(3, 3, nil)
	rn.Sub(t.rn_s, mean_multi)

	// Converting Rn to a Symmetric Matrix format
	backing := make([]float64, 0, 9)
	for i := 0; i < 3; i++ {
		for j := range rn.RawRowView(i) {
			backing = append(backing, rn.RawRowView(i)[j])
		}
	}

	t.rn = mat.NewSymDense(3, backing)
}

// Calculates the largest eigenvector of the node, returns the vector and its value
func (t *tree_node) _calculate_max_eigenvector() (*mat.VecDense, float64) {
	// NaN indicates the node has zero pixels so the eigen value cannot be calculated
	if math.IsNaN(t.rn.At(0, 0)) {
		return mat.NewVecDense(3, nil), 0
	}

	// Calculating the eigenvalues
	var eigsym mat.EigenSym
	ok := eigsym.Factorize(t.rn, true)
	chkFatal(ok, "factorizing eigen not worked")
	//fmt.Printf("Eigenvalues of:\n%1.3f\n\n", eigsym.Values(nil))

	// Calculating the eigenvector
	var ev mat.Dense
	eigsym.VectorsTo(&ev)
	//fmt.Printf("Eigenvectors of:\n%1.3f\n\n", mat.Formatted(&ev))

	// Finding the largest eigenvector for the current node
	var current_vector *mat.VecDense
	var current_vector_magnitude float64

	row_0_vec := mat.VecDenseCopyOf(ev.RowView(0))
	row_1_vec := mat.VecDenseCopyOf(ev.RowView(1))
	row_2_vec := mat.VecDenseCopyOf(ev.RowView(2))
	row_0_vec_magnitude := _3DVecMagnitude(row_0_vec)
	row_1_vec_magnitude := _3DVecMagnitude(row_1_vec)
	row_2_vec_magnitude := _3DVecMagnitude(row_2_vec)

	current_vector = row_0_vec
	current_vector_magnitude = row_0_vec_magnitude
	if row_1_vec_magnitude > row_0_vec_magnitude {
		current_vector = row_1_vec
		current_vector_magnitude = row_1_vec_magnitude
	} else if row_2_vec_magnitude > row_1_vec_magnitude {
		current_vector = row_2_vec
		current_vector_magnitude = row_2_vec_magnitude
	}

	//if current_vector.At(0, 0) == 1 || current_vector.At(1, 0) == 1 || current_vector.At(2, 0) == 1 {
	//	return current_vector, 0
	//}

	return current_vector, current_vector_magnitude
}

// Called on the root node, splits the node in the root's tree with the largest eigenvector
func (t *tree_node) partition_node() {
	set1 := gs.NewSet()
	set2 := gs.NewSet()

	eigenvec, _ := t._calculate_max_eigenvector()

	oneonematrix := mat.NewDense(1, 1, nil)
	//matPrint(eigenvec.T())
	oneonematrix.Mul(eigenvec.T(), t.qn)
	limit_val := oneonematrix.At(0, 0)

	for _, v := range t.pixels.List() {
		matrix, ok := v.(mat.Matrix)
		chkFatal(ok, "retrieving pixels when partitioning node not worked")

		oneonematrix.Mul(eigenvec.T(), matrix)
		value := oneonematrix.At(0, 0)

		if value <= limit_val {
			set1.Add(matrix)
			t.pixels.Remove(matrix)
		} else {
			set2.Add(matrix)
			t.pixels.Remove(matrix)
		}

	}

	t.left = newNode(set1)
	t.right = newNode(set2)

	t.left.calculate_mean_and_covariance()
	//t.right.calculate_mean_and_covariance()
	t.applyRelation()
}

// Instead of fully calculating the mean and covariance of the right node, the relation
// is applied to make the program more efficient
func (t *tree_node) applyRelation() {
	// Calculating rn non squiggly
	t.right.rn_s.Sub(t.rn_s, t.left.rn_s)

	// Calculating mn
	t.right.mn.Sub(t.mn, t.left.mn)

	// Calculating nn
	t.right.nn = 0
	t.right.nn = t.nn - t.left.nn

	// Calculates Qn by dividing Mn by the Cardinality of the Set
	t.right.qn.Scale(float64(1)/float64(t.right.nn), t.right.mn)

	// Calculating rn squiggly
	var mean_multi = mat.NewDense(3, 3, nil)
	mean_multi.Mul(t.right.qn, t.right.qn.T())

	// Minus the mean from Rn (non-squiggle) to get the covariance (R squiggly)
	rn := mat.NewDense(3, 3, nil)
	rn.Sub(t.right.rn_s, mean_multi)

	// Converting Rn to a Symmetric Matrix format
	backing := make([]float64, 0, 9)
	for i := 0; i < 3; i++ {
		for j := range rn.RawRowView(i) {
			backing = append(backing, rn.RawRowView(i)[j])
		}
	}

	t.right.rn = mat.NewSymDense(3, backing)

}

// Finds the node with the largest eigenvector in the root tree (call on the root)
func (t *tree_node) find_max_eigenvector() *tree_node {
	q := golang_queues.New()
	q.Enqueue(t)

	var max_value float64
	var max_node *tree_node
	for q.Len() > 0 {
		front := q.Front().Value
		node, ok := front.(*tree_node)
		chkFatal(ok, "retrieving front of queue when finding max eigen vector not working")

		var val float64
		if node.left != nil || node.right != nil {
			// Ignored node if it has a child meaning it has been visited
			val = 0

			if node.left != nil {
				q.Enqueue(node.left)
			}
			if node.right != nil {
				q.Enqueue(node.right)
			}

		} else {
			// If node not visited then it calculates the eigenvector value
			_, val = node._calculate_max_eigenvector()
		}

		if val > max_value {
			max_value = val
			max_node = node
		}

		q.Dequeue()
	}

	return max_node
}

// Returns all leaaves in the nodes tree
func (t *tree_node) get_leaves() []*tree_node {
	q := golang_queues.New()
	q.Enqueue(t)

	leaves := make([]*tree_node, 0, 5)

	for q.Len() > 0 {
		front := q.Front().Value
		node, ok := front.(*tree_node)
		chkFatal(ok, "retrieving front of queue when finding leaves not working")

		if node.left != nil {
			q.Enqueue(node.left)
		}
		if node.right != nil {
			q.Enqueue(node.right)
		}
		if node.left == nil && node.right == nil {
			leaves = append(leaves, node)
		}

		q.Dequeue()
	}

	return leaves[:len(leaves)-1]
}

// Find "n" most dominant colours with the path to the image given
// Uses a custom binary tree
func FindDominantColoursBasic(path string, n int) []*RGB {
	root := newNode(Img2pixelset(path))
	root.calculate_mean_and_covariance()

	for i := 0; i < n; i++ {
		fmt.Printf("WORKING ON ITERATION: %v\n", i)
		node := root.find_max_eigenvector()
		node.calculate_mean_and_covariance()
		node.partition_node()
	}

	colours := make([]*RGB, 0, n)
	leaves := root.get_leaves()
	for i := range leaves {
		r := leaves[i].qn.At(0, 0)
		g := leaves[i].qn.At(1, 0)
		b := leaves[i].qn.At(2, 0)
		colours = append(colours, &RGB{r, g, b})
	}

	return colours
}

// If you want to reuse the same pixel set multiple times you can read
// the data into a set using the Img3pixelset function and then pass it
// into this function which makes a copy of the set before calculating
// the dominant colours leaving the original set untouched
func FindDominantColoursBasicFromSet(s *gs.Set, n int) []*RGB {
	root := newNode(s)
	root.calculate_mean_and_covariance()

	for i := 0; i < n; i++ {
		fmt.Printf("WORKING ON ITERATION: %v\n", i)
		node := root.find_max_eigenvector()
		node.calculate_mean_and_covariance()
		node.partition_node()
	}

	colours := make([]*RGB, 0, n)
	leaves := root.get_leaves()
	for i := range leaves {
		r := leaves[i].qn.At(0, 0)
		g := leaves[i].qn.At(1, 0)
		b := leaves[i].qn.At(2, 0)
		colours = append(colours, &RGB{r, g, b})
	}

	return colours
}
