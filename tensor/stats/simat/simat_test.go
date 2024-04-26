// Copyright (c) 2024, Cogent Core. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simat

import (
	"testing"

	"cogentcore.org/core/tensor/stats/metric"
	"cogentcore.org/core/tensor/table"

	"github.com/stretchr/testify/assert"
)

var simres = `Tensor: [12, 12]
[0]:       0 3.4641016151377544 8.831760866327848 9.273618495495704 8.717797887081348 9.38083151964686 4.69041575982343 5.830951894845301 8.12403840463596 8.54400374531753 5.291502622129181 6.324555320336759 
[1]: 3.4641016151377544       0 9.38083151964686 8.717797887081348 9.273618495495704 8.831760866327848 5.830951894845301 4.69041575982343 8.717797887081348 7.937253933193772 6.324555320336759 5.291502622129181 
[2]: 8.831760866327848 9.38083151964686       0 3.4641016151377544 4.242640687119285 5.0990195135927845 9.38083151964686 9.899494936611665 4.47213595499958 5.744562646538029 9.38083151964686 9.899494936611665 
[3]: 9.273618495495704 8.717797887081348 3.4641016151377544       0 5.477225575051661 3.7416573867739413 9.797958971132712 9.273618495495704 5.656854249492381 4.58257569495584 9.797958971132712 9.273618495495704 
[4]: 8.717797887081348 9.273618495495704 4.242640687119285 5.477225575051661       0       4 8.831760866327848 9.38083151964686 4.242640687119285 5.5677643628300215 8.831760866327848 9.38083151964686 
[5]: 9.38083151964686 8.831760866327848 5.0990195135927845 3.7416573867739413       4       0 9.486832980505138 8.94427190999916 5.830951894845301 4.795831523312719 9.486832980505138 8.94427190999916 
[6]: 4.69041575982343 5.830951894845301 9.38083151964686 9.797958971132712 8.831760866327848 9.486832980505138       0 3.4641016151377544 9.16515138991168 9.539392014169456 4.242640687119285 5.477225575051661 
[7]: 5.830951894845301 4.69041575982343 9.899494936611665 9.273618495495704 9.38083151964686 8.94427190999916 3.4641016151377544       0 9.695359714832659       9 5.477225575051661 4.242640687119285 
[8]: 8.12403840463596 8.717797887081348 4.47213595499958 5.656854249492381 4.242640687119285 5.830951894845301 9.16515138991168 9.695359714832659       0 3.605551275463989 9.16515138991168 9.695359714832659 
[9]: 8.54400374531753 7.937253933193772 5.744562646538029 4.58257569495584 5.5677643628300215 4.795831523312719 9.539392014169456       9 3.605551275463989       0 9.539392014169456       9 
[10]: 5.291502622129181 6.324555320336759 9.38083151964686 9.797958971132712 8.831760866327848 9.486832980505138 4.242640687119285 5.477225575051661 9.16515138991168 9.539392014169456       0 3.4641016151377544 
[11]: 6.324555320336759 5.291502622129181 9.899494936611665 9.273618495495704 9.38083151964686 8.94427190999916 5.477225575051661 4.242640687119285 9.695359714832659       9 3.4641016151377544       0 
`

func TestClust(t *testing.T) {
	dt := &table.Table{}
	err := dt.OpenCSV("../clust/testdata/faces.dat", table.Tab)
	if err != nil {
		t.Error(err)
	}
	ix := table.NewIndexView(dt)
	smat := &SimMat{}
	smat.TableCol(ix, "Input", "Name", false, metric.Euclidean64)

	// fmt.Println(smat.Mat)
	assert.Equal(t, simres, smat.Mat.String())
}