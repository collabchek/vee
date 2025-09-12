package vee

import (
	"testing"
	"time"
)

func TestTimeRendering(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{
			name: "time.Time field renders as datetime-local by default",
			input: struct {
				Created time.Time
			}{Created: time.Date(2023, 12, 25, 14, 30, 0, 0, time.UTC)},
			want: `<form method="POST">
<label for="created">Created</label>
<input type="datetime-local" name="created" value="2023-12-25T14:30" id="created">
</form>
`,
		},
		{
			name: "time.Time with zero value renders without value",
			input: struct {
				Created time.Time
			}{Created: time.Time{}},
			want: `<form method="POST">
<label for="created">Created</label>
<input type="datetime-local" name="created" id="created">
</form>
`,
		},
		{
			name: "time.Time with type='date' renders as date input",
			input: struct {
				Birthday time.Time `vee:"type:'date'"`
			}{Birthday: time.Date(1990, 6, 15, 0, 0, 0, 0, time.UTC)},
			want: `<form method="POST">
<label for="birthday">Birthday</label>
<input type="date" name="birthday" value="1990-06-15" id="birthday">
</form>
`,
		},
		{
			name: "time.Time with type='time' renders as time input",
			input: struct {
				MeetingTime time.Time `vee:"type:'time'"`
			}{MeetingTime: time.Date(2023, 1, 1, 15, 45, 0, 0, time.UTC)},
			want: `<form method="POST">
<label for="meeting_time">Meeting Time</label>
<input type="time" name="meeting_time" value="15:45" id="meeting_time">
</form>
`,
		},
		{
			name: "time.Time with min/max attributes",
			input: struct {
				Deadline time.Time `vee:"type:'date',min:'2023-01-01',max:'2023-12-31'"`
			}{Deadline: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
			want: `<form method="POST">
<label for="deadline">Deadline</label>
<input type="date" name="deadline" value="2023-06-15" min="2023-01-01" max="2023-12-31" id="deadline">
</form>
`,
		},
		{
			name: "time.Time with custom name override",
			input: struct {
				StartTime time.Time `vee:"$start_at,type:'datetime-local'"`
			}{StartTime: time.Date(2023, 8, 10, 9, 0, 0, 0, time.UTC)},
			want: `<form method="POST">
<label for="start_at">Start Time</label>
<input type="datetime-local" name="start_at" value="2023-08-10T09:00" id="start_at">
</form>
`,
		},
		{
			name: "time.Duration field renders as number input with default seconds",
			input: struct {
				Timeout time.Duration
			}{Timeout: 30 * time.Second},
			want: `<form method="POST">
<label for="timeout">Timeout</label>
<input type="number" name="timeout" value="30" id="timeout">
</form>
`,
		},
		{
			name: "time.Duration with zero value renders without value",
			input: struct {
				Timeout time.Duration
			}{Timeout: 0},
			want: `<form method="POST">
<label for="timeout">Timeout</label>
<input type="number" name="timeout" id="timeout">
</form>
`,
		},
		{
			name: "time.Duration with units='ms'",
			input: struct {
				Delay time.Duration `vee:"units:'ms'"`
			}{Delay: 500 * time.Millisecond},
			want: `<form method="POST">
<label for="delay">Delay</label>
<input type="number" name="delay" value="500" id="delay">
</form>
`,
		},
		{
			name: "time.Duration with units='m' (minutes)",
			input: struct {
				Duration time.Duration `vee:"units:'m'"`
			}{Duration: 2*time.Hour + 30*time.Minute},
			want: `<form method="POST">
<label for="duration">Duration</label>
<input type="number" name="duration" value="150" id="duration">
</form>
`,
		},
		{
			name: "time.Duration with units='h' (hours)",
			input: struct {
				WorkDay time.Duration `vee:"units:'h'"`
			}{WorkDay: 8 * time.Hour},
			want: `<form method="POST">
<label for="work_day">Work Day</label>
<input type="number" name="work_day" value="8" id="work_day">
</form>
`,
		},
		{
			name: "time.Duration with min/max/step attributes",
			input: struct {
				Timeout time.Duration `vee:"units:'s',min:'1',max:'3600',step:'1'"`
			}{Timeout: 60 * time.Second},
			want: `<form method="POST">
<label for="timeout">Timeout</label>
<input type="number" name="timeout" value="60" min="1" max="3600" step="1" id="timeout">
</form>
`,
		},
		{
			name: "time.Duration with custom name override",
			input: struct {
				MaxWait time.Duration `vee:"$max_wait_time,units:'ms'"`
			}{MaxWait: 250 * time.Millisecond},
			want: `<form method="POST">
<label for="max_wait_time">Max Wait</label>
<input type="number" name="max_wait_time" value="250" id="max_wait_time">
</form>
`,
		},
		{
			name: "mixed types with time fields",
			input: struct {
				Name      string
				CreatedAt time.Time     `vee:"type:'datetime-local'"`
				ExpiresIn time.Duration `vee:"units:'h'"`
				Active    bool
			}{
				Name:      "Test Event",
				CreatedAt: time.Date(2023, 12, 1, 10, 0, 0, 0, time.UTC),
				ExpiresIn: 24 * time.Hour,
				Active:    true,
			},
			want: `<form method="POST">
<label for="name">Name</label>
<input type="text" name="name" value="Test Event" id="name">
<label for="created_at">Created At</label>
<input type="datetime-local" name="created_at" value="2023-12-01T10:00" id="created_at">
<label for="expires_in">Expires In</label>
<input type="number" name="expires_in" value="24" id="expires_in">
<label for="active">Active</label>
<input type="checkbox" name="active" value="true" checked id="active">
</form>
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.input)
			if err != nil {
				t.Errorf("Render() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("Render() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTimeBinding(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string][]string
		target  func() any
		check   func(t *testing.T, target any)
		wantErr bool
	}{
		{
			name: "time.Time datetime-local binding",
			input: map[string][]string{
				"created": {"2023-12-25T14:30"},
			},
			target: func() any {
				return &struct {
					Created time.Time
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct{ Created time.Time })
				expected := time.Date(2023, 12, 25, 14, 30, 0, 0, time.UTC)
				if !s.Created.Equal(expected) {
					t.Errorf("Expected Created=%v, got Created=%v", expected, s.Created)
				}
			},
		},
		{
			name: "time.Time date binding",
			input: map[string][]string{
				"birthday": {"1990-06-15"},
			},
			target: func() any {
				return &struct {
					Birthday time.Time `vee:"type:'date'"`
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					Birthday time.Time `vee:"type:'date'"`
				})
				expected := time.Date(1990, 6, 15, 0, 0, 0, 0, time.UTC)
				if !s.Birthday.Equal(expected) {
					t.Errorf("Expected Birthday=%v, got Birthday=%v", expected, s.Birthday)
				}
			},
		},
		{
			name: "time.Time time binding",
			input: map[string][]string{
				"meeting_time": {"15:45"},
			},
			target: func() any {
				return &struct {
					MeetingTime time.Time `vee:"type:'time'"`
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					MeetingTime time.Time `vee:"type:'time'"`
				})
				expected := time.Date(0, 1, 1, 15, 45, 0, 0, time.UTC)
				if !s.MeetingTime.Equal(expected) {
					t.Errorf("Expected MeetingTime=%v, got MeetingTime=%v", expected, s.MeetingTime)
				}
			},
		},
		{
			name: "time.Time with custom name override",
			input: map[string][]string{
				"start_at": {"2023-08-10T09:00"},
			},
			target: func() any {
				return &struct {
					StartTime time.Time `vee:"$start_at,type:'datetime-local'"`
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					StartTime time.Time `vee:"$start_at,type:'datetime-local'"`
				})
				expected := time.Date(2023, 8, 10, 9, 0, 0, 0, time.UTC)
				if !s.StartTime.Equal(expected) {
					t.Errorf("Expected StartTime=%v, got StartTime=%v", expected, s.StartTime)
				}
			},
		},
		{
			name: "time.Duration seconds binding (default)",
			input: map[string][]string{
				"timeout": {"30"},
			},
			target: func() any {
				return &struct {
					Timeout time.Duration
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct{ Timeout time.Duration })
				expected := 30 * time.Second
				if s.Timeout != expected {
					t.Errorf("Expected Timeout=%v, got Timeout=%v", expected, s.Timeout)
				}
			},
		},
		{
			name: "time.Duration milliseconds binding",
			input: map[string][]string{
				"delay": {"500"},
			},
			target: func() any {
				return &struct {
					Delay time.Duration `vee:"units:'ms'"`
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					Delay time.Duration `vee:"units:'ms'"`
				})
				expected := 500 * time.Millisecond
				if s.Delay != expected {
					t.Errorf("Expected Delay=%v, got Delay=%v", expected, s.Delay)
				}
			},
		},
		{
			name: "time.Duration minutes binding",
			input: map[string][]string{
				"duration": {"150"},
			},
			target: func() any {
				return &struct {
					Duration time.Duration `vee:"units:'m'"`
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					Duration time.Duration `vee:"units:'m'"`
				})
				expected := 150 * time.Minute
				if s.Duration != expected {
					t.Errorf("Expected Duration=%v, got Duration=%v", expected, s.Duration)
				}
			},
		},
		{
			name: "time.Duration hours binding",
			input: map[string][]string{
				"work_day": {"8"},
			},
			target: func() any {
				return &struct {
					WorkDay time.Duration `vee:"units:'h'"`
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					WorkDay time.Duration `vee:"units:'h'"`
				})
				expected := 8 * time.Hour
				if s.WorkDay != expected {
					t.Errorf("Expected WorkDay=%v, got WorkDay=%v", expected, s.WorkDay)
				}
			},
		},
		{
			name: "time.Duration with custom name override",
			input: map[string][]string{
				"max_wait_time": {"250"},
			},
			target: func() any {
				return &struct {
					MaxWait time.Duration `vee:"$max_wait_time,units:'ms'"`
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					MaxWait time.Duration `vee:"$max_wait_time,units:'ms'"`
				})
				expected := 250 * time.Millisecond
				if s.MaxWait != expected {
					t.Errorf("Expected MaxWait=%v, got MaxWait=%v", expected, s.MaxWait)
				}
			},
		},
		{
			name: "mixed fields with time types",
			input: map[string][]string{
				"name":       {"Test Event"},
				"created_at": {"2023-12-01T10:00"},
				"expires_in": {"24"},
				"active":     {"true"},
			},
			target: func() any {
				return &struct {
					Name      string
					CreatedAt time.Time     `vee:"type:'datetime-local'"`
					ExpiresIn time.Duration `vee:"units:'h'"`
					Active    bool
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					Name      string
					CreatedAt time.Time     `vee:"type:'datetime-local'"`
					ExpiresIn time.Duration `vee:"units:'h'"`
					Active    bool
				})
				expectedTime := time.Date(2023, 12, 1, 10, 0, 0, 0, time.UTC)
				expectedDuration := 24 * time.Hour

				if s.Name != "Test Event" || !s.CreatedAt.Equal(expectedTime) || s.ExpiresIn != expectedDuration || !s.Active {
					t.Errorf("Expected Name='Test Event' CreatedAt=%v ExpiresIn=%v Active=true, got Name='%s' CreatedAt=%v ExpiresIn=%v Active=%t",
						expectedTime, expectedDuration, s.Name, s.CreatedAt, s.ExpiresIn, s.Active)
				}
			},
		},
		{
			name: "empty time fields don't bind",
			input: map[string][]string{
				"name": {"Test"},
				// created and timeout absent
			},
			target: func() any {
				return &struct {
					Name    string
					Created time.Time
					Timeout time.Duration
				}{
					// Set initial values to non-zero to ensure they don't change
					Created: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					Timeout: 5 * time.Second,
				}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					Name    string
					Created time.Time
					Timeout time.Duration
				})
				expectedTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
				expectedDuration := 5 * time.Second

				if s.Name != "Test" || !s.Created.Equal(expectedTime) || s.Timeout != expectedDuration {
					t.Errorf("Expected Name='Test' Created=%v Timeout=%v (unchanged), got Name='%s' Created=%v Timeout=%v",
						expectedTime, expectedDuration, s.Name, s.Created, s.Timeout)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := tt.target()
			err := Bind(tt.input, target)

			if (err != nil) != tt.wantErr {
				t.Errorf("Bind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				tt.check(t, target)
			}
		})
	}
}
