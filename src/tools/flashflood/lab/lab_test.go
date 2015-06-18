package lab_test

import (
	"tools/flashflood/lab"
	"tools/flashflood/experiment"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lab", func() {
	It("starts an experiment", func() {
		fe := fakeExperiment{}

		var newFE = func() lab.Experiment {
			return &fe
		}

		l := lab.New(newFE)
		l.Start()
		Eventually(fe.Started).Should(BeTrue())
	})

	It("stops a running experiment", func() {
		fe := fakeExperiment{}

		var newFE = func() lab.Experiment {
			return &fe
		}

		l := lab.New(newFE)
		l.Start()
		Eventually(fe.Started).Should(BeTrue())
		l.Stop()
		Expect(fe.stop).To(BeClosed())
	})

	It("starts a new experiment if start called while one is running", func() {
		fes := []*fakeExperiment{}

		var newFE = func() lab.Experiment {
			fe := fakeExperiment{}
			fes = append(fes, &fe)
			return &fe
		}

		l := lab.New(newFE)

		l.Start()
		Eventually(fes[0].Started).Should(BeTrue())
		l.Start()
		Eventually(fes[1].Started).Should(BeTrue())

		Expect(fes).To(HaveLen(2))
		Expect(fes[0].stop).To(BeClosed())
		Expect(fes[1].stop).NotTo(BeClosed())
	})

	It("reports results of latest experiment", func() {
		fe := fakeExperiment{}

		var newFE = func() lab.Experiment {
			return &fe
		}

		l := lab.New(newFE)
		l.Start()

		l.Results()
		Eventually(fe.ResultsCalled).Should(BeTrue())
	})
})

type fakeExperiment struct {
	stop <-chan struct{}

	started bool
	resultsCalled bool
}

func (f *fakeExperiment) Start(stop <-chan struct{}) {
	f.stop = stop
	f.started = true
}

func (f *fakeExperiment) Results() experiment.Results {
	f.resultsCalled = true
	return []experiment.Result{}
}

func (f *fakeExperiment) ResultsCalled() bool {
	return f.resultsCalled
}

func (f *fakeExperiment) Started() bool {
	return f.started
}
