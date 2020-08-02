package controller

/*func TestCreateProbeHandler(t *testing.T) {
	// Insert the JSON body as a string.
	jsonBody := `{
      "Name": "simple-service-http",
      "URL": "http://localhost:8080/actuator/health",
      "Delay": 5
    }`

	// Insert a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(http.MethodPost, "/api/v1/probe/create", strings.NewReader(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	controller := NewProbeController(&fakeProbeService{})

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.Create)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `CreateProbeRequest: {Name:simple-service-http URL:http://localhost:8080/actuator/health Delay:5}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

type fakeProbeService struct{}

func (ps *fakeProbeService) Create(_ prober.Probe) error {
	fmt.Println("Calling Insert() from mock ProbeService.")
	return nil
}*/
