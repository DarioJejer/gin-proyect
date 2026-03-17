# Controller Test Specifications – Prompts for Agent
---

## Style and conventions

- Use `testify/suite` with `UnitTestSuite` and `IntegrationTestSuite` structs.
- Test names: `Test_<HTTP_Method>_<Scenario>_<ExpectedStatus>` (e.g. `Test_Post_ValidCreation_StatusCreated`).
- Unit tests: inject `mocks.IUsersRepository`, `mocks.ICompaniesRepository`; use `EXPECT()` to set return values.
- Integration tests: use `initializers.LoadEnvVariables`, `initializers.ConnectToDB`, real `repositories.NewUsersRepository()` / `repositories.NewCompaniesRepository()`.
- Integration tests use `SetupSuite` (DB connect, controller init) and `SetupTest` (delete tables content).
- An every test always `gin.SetMode(gin.TestMode)`, create `router := gin.Default()`, register the handler, use `httptest.NewRecorder()` and `http.NewRequest`.
- Assert status with `assert.Equal(suite.T(), expectedStatus, r.Code)`.

### Response validation (success responses)

For tests that expect a successful response (2xx), also validate the response body:

- **Single resource (GET by id, POST create, PUT update):** Parse the JSON response and assert that key fields match the expected values (e.g. `user.name`, `user.age`, `user.company_id`, `user.id`; or `company.name`, `company.id`). Use the same field names as in the API (e.g. `user`, `company` wrapper keys).
- **List (GET collection):** Parse the response and assert the **number of items** (e.g. `len(users) == 0` for empty, `len(users) == 2` when two were created). Where the test creates specific data, also assert that at least one key field (e.g. name) matches in the returned items so the list content is correct.

Use `json.Unmarshal(r.Body.Bytes(), &target)` or a decoder with a struct that matches the API response shape (e.g. `gin.H`-style wrappers like `{"user": ..., "status": ...}` or `{"users": [...]}`). Prefer type-safe structs (e.g. `responseDTOs.UserResponseDTO`) for the nested payload.

---

## UsersController

### Unit tests (mocks, no DB)

| **Test_Post_ValidCreation_StatusCreated** – POST /users with valid `CreateUserDTO` (name, age, company_id). Mock `CompaniesRepo.GetCompany` to return a company, mock `UserRepo.PostUser` to return nil. Expect 201. Assert response body: `status` "user created", `user` object with matching `name`, `age`, `company_id`, and `id` > 0. |
| **Test_Post_InvalidUserData_StatusBadRequest** – POST /users with invalid data (e.g. missing age or company_id). No mocks needed. Expect 400. |
| **Test_Post_InvalidBody_StatusBadRequest** – POST /users with malformed JSON (e.g. empty body or invalid syntax). Expect 400. |
| **Test_Post_CompanyNotFound_StatusBadRequest** – POST /users with valid DTO but mock `CompaniesRepo.GetCompany` to return error. Expect 400. |
| **Test_Post_RepoError_StatusInternalServerError** – POST /users with valid DTO, mock GetCompany to succeed, mock `UserRepo.PostUser` to return error. Expect 500. |
| **Test_Get_ValidId_StatusOk** – GET /users/:id with valid ID. Mock `UserRepo.GetUser` to return a user. Expect 200. Assert response body: `user` object with correct `id`, `name`, `age`, `company_id` matching the mocked user. |
| **Test_Get_InvalidId_StatusBadRequest** – GET /users/:id with non-numeric ID (e.g. "abc"). Expect 400. |
| **Test_Get_UserNotFound_StatusNotFound** – GET /users/:id with valid numeric ID. Mock `UserRepo.GetUser` to return `gorm.ErrRecordNotFound`. Expect 404. |
| **Test_GetUsers_Valid_StatusOk** – GET /users. Mock `UserRepo.GetUsers` to return a slice of users. Expect 200. Assert response body: `users` array length matches mocked slice (e.g. 1); first item has correct `name` (and optionally `id`, `age`, `company_id`). |
| **Test_GetUsers_RepoError_StatusInternalServerError** – GET /users. Mock `UserRepo.GetUsers` to return error. Expect 500. |
| **Test_Update_ValidCreation_StatusOk** – PUT/PATCH /users with valid DTO. Mock `CompaniesRepo.GetCompany` and `UserRepo.UpdateUser`. Expect 200. Assert response body: `status` "user updated", `user` object with updated `name`, `age`, `company_id` matching the request DTO. |
| **Test_Update_InvalidUserData_StatusBadRequest** – PUT/PATCH /users with invalid DTO (missing required fields). Expect 400. |
| **Test_Update_InvalidBody_StatusBadRequest** – PUT /users with malformed JSON (e.g. invalid syntax). Expect 400. |
| **Test_Update_CompanyNotFound_StatusBadRequest** – PUT /users with valid DTO but mock `CompaniesRepo.GetCompany` to return error. Expect 400. |
| **Test_Update_RepoError_StatusInternalServerError** – PUT /users with valid DTO. Mock `UserRepo.UpdateUser` to return error. Expect 500. |

### Integration tests (real DB)

| **Test_Post_ValidCreation_StatusCreated** – POST /users Setup: create company in DB. Create valid user using company.ID. Expect 201. Assert response: `status` "user created", `user` with matching `name`, `age`, `company_id`, and `id` > 0. |
| **Test_Post_InvalidCompany_StatusBadRequest** – POST /users with non-existent company_id. Expect 400. |
| **Test_Post_InvalidUserData_StatusBadRequest** – POST /users with invalid DTO (e.g. missing age). Expect 400. |
| **Test_Post_InvalidBody_StatusBadRequest** – POST /users with malformed JSON (e.g. invalid syntax). Expect 400. |
| **Test_Get_ValidId_StatusOk** – GET /users/:id. Setup: Create company and then user in DB using company.ID. Then GET by user.ID. Expect 200. Assert response: `user` with same `id`, `name`, `age`, `company_id` as created. |
| **Test_Get_InvalidId_StatusBadRequest** – GET /users/:id with non-numeric ID (e.g. "abc"). Expect 400. |
| **Test_Get_UserNotFound_StatusNotFound** – GET /users/:id with ID that does not exist (e.g. "9999"). Expect 404. |
| **Test_GetUsers_Empty_StatusOk** – GET /users with empty DB. Expect 200. Assert response: `users` array with length 0. |
| **Test_GetUsers_WithData_StatusOk** – GET /users. Setup: create company and two users using company.ID in DB. Then GET users. Expect 200. Assert response: `users` array length 2; assert names (or ids) match the two created users. |
| **Test_Update_Valid_StatusOk** – PUT /users/:id. Setup: Create user and company in DB. Then PUT with valid DTO and existing company. Expect 200. Assert response: `status` "user updated", `user` with updated `name`, `age`, `company_id` matching the request. |
| **Test_Update_InvalidBody_StatusBadRequest** – PUT /users/:id with malformed JSON (e.g. invalid syntax). Expect 400. |
| **Test_Update_InvalidUserData_StatusBadRequest** – PUT /users/:id with invalid DTO (e.g. missing age). Expect 400. |
| **Test_Update_InvalidId_StatusBadRequest** – PUT /users/:id with valid DTO but non-numeric ID (e.g. "abc"). Expect 400. |
| **Test_Update_UserNotFound_StatusNotFound** – PUT /users/99999 with valid DTO. Expect 404. |
| **Test_Update_InvalidCompany_StatusBadRequest** – PUT /users/:id with valid DTO and non-existent company_id. Expect 400. |


---

## CompaniesController

### Unit tests (mocks, no DB)

| **Test_Post_ValidCreation_StatusCreated** – POST /companies with valid `CreateCompanyDTO` (name). Mock `CompanyRepo.PostCompany` to return nil. Expect 201. Assert response: `status` "company created", `company` with matching `name` and `id` > 0. |
| **Test_Post_InvalidCompanyData_StatusBadRequest** – POST /companies with invalid data (e.g. empty name or missing name). No mocks needed. Expect 400. |
| **Test_Post_InvalidBody_StatusBadRequest** – POST /companies with malformed JSON (e.g. empty body or invalid syntax). Expect 400. |
| **Test_Post_RepoError_StatusInternalServerError** – POST /companies with valid DTO. Mock `CompanyRepo.PostCompany` to return error. Expect 500. |
| **Test_Get_ValidId_StatusOk** – GET /companies/:id with valid ID. Mock `CompanyRepo.GetCompany` to return a company. Expect 200. Assert response: `company` with correct `id`, `name` matching mocked company. |
| **Test_Get_InvalidId_StatusBadRequest** – GET /companies/:id with non-numeric ID (e.g. "abc"). Expect 400. |
| **Test_Get_CompanyNotFound_StatusNotFound** – GET /companies/:id with valid numeric ID. Mock `CompanyRepo.GetCompany` to return `gorm.ErrRecordNotFound`. Expect 404. |
| **Test_Get_RepoError_StatusInternalServerError** – GET /companies/:id with valid ID. Mock `CompanyRepo.GetCompany` to return error. Expect 500. |
| **Test_GetCompanies_Valid_StatusOk** – GET /companies. Mock `CompanyRepo.GetCompanies` to return a slice of companies. Expect 200. Assert response: `companies` array length matches mocked slice (e.g. 1); first item has correct `name` and `id`. |
| **Test_GetCompanies_RepoError_StatusInternalServerError** – GET /companies. Mock `CompanyRepo.GetCompanies` to return error. Expect 500. |

### Integration tests (real DB)

| **Test_Post_ValidCreation_StatusCreated** – POST /companies. Setup: none. Send valid DTO. Expect 201. Assert response: `status` "company created", `company` with matching `name` and `id` > 0. |
| **Test_Post_InvalidCompanyData_StatusBadRequest** – POST /companies with invalid DTO (e.g. empty name). Expect 400. |
| **Test_Post_InvalidBody_StatusBadRequest** – POST /companies with malformed JSON (e.g. invalid syntax). Expect 400. |
| **Test_Get_ValidId_StatusOk** – GET /companies/:id. Setup: create company in DB. Then GET by company.ID. Expect 200. Assert response: `company` with same `id`, `name` as created. |
| **Test_Get_InvalidId_StatusBadRequest** – GET /companies/:id with non-numeric ID (e.g. "abc"). Expect 400. |
| **Test_Get_CompanyNotFound_StatusNotFound** – GET /companies/:id with ID that does not exist (e.g. "99999"). Expect 404. |
| **Test_GetCompanies_Empty_StatusOk** – GET /companies with empty DB. Expect 200. Assert response: `companies` array with length 0. |
| **Test_GetCompanies_WithData_StatusOk** – GET /companies. Setup: create two companies in DB. Then GET /companies. Expect 200. Assert response: `companies` array length 2; assert names (or ids) match the two created companies. |

---

## File layout

- Unit tests: create `controllers/usersController_unit_test.go` for UsersController.
- Integration tests: create `controllers/usersController_integration_test.go` for UsersController.

---

## DTOs reference

- `CreateUserDTO`: `Name` (string, required), `Age` (int, required, gt=0), `CompanyID` (uint, required, gt=0).
- `CreateCompanyDTO`: `Name` (string, required).
- `UserResponseDTO`: `ID` (uint), `Name` (string), `Age` (int), `CompanyID` (uint)
- `CompanyResponseDTO`: `ID` (uint), `Name` (string)
