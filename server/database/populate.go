package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// shot needs name description and id

func PopulateDB(ctx context.Context, pool *pgxpool.Pool) {
	makeAnimals(ctx, pool)
	populateShots(ctx, pool)
	makeImage(ctx, pool)
	populateNotificationTypes(pool)
}

func makeAnimals(ctx context.Context, pool *pgxpool.Pool) {
	sql := []string{
		` INSERT INTO animals (species, date_of_birth, sex, available, breed, name) VALUES (
			'dog', '2022/10/14', 'male', true, 'Bichon Frise Poodle', 'Bimbus');`,
		` INSERT INTO animals (species, date_of_birth, sex, price, available, breed, name) VALUES (
			'cat', '2020/09/11', 'undefined' , 9.11, false, 'Russian Blue', 'Florida');`,
		` INSERT INTO animals (species, date_of_birth, sex, price, available, breed, name) VALUES (
			'cat', '8008/09/11', 'female', 100.00, true, 'calico', 'Garu');`,
	}
	for _, query := range sql {

		_, err := pool.Exec(ctx, query)
		queryFail(err, "insert animal data")
	}
}

func makeImage(ctx context.Context, pool *pgxpool.Pool) {
	sql := []string{
		`INSERT INTO animal_images (animal_id, image_src, priority) VALUES ('90f74cb5-2f24-434d-90ac-4553c9f31cd4', './images/dog.webp', 1);`,
	}

	for _, query := range sql {
		_, err := pool.Exec(ctx, query)
		queryFail(err, "insert image data")
	}
}

func populateShots(ctx context.Context, pool *pgxpool.Pool) {
	sql := []string{
		`INSERT INTO shots (name, description) VALUES ('rabies','A vital vaccine that builds immunity against a deadly viral disease that affects the central nervous system.');`,
		`INSERT INTO shots(name,description) VALUES('distemper','Protects against a debilitating virus that affects the respiratory, gastrointestinal, and nervous systems');`,
		`INSERT INTO shots(name,description) VALUES('lyme disease','Caused by a bacterium transmitted through tick bites, this vaccine can help protect dogs that are frequently exposed to ticks');`,
		`INSERT INTO shots(name,description) VALUES('bordetella bronchiseptica','Also known as "Kennel Cough", this vaccine is recommended for dogs that interact with large groups of unknown dogs');`,
		`INSERT INTO shots(name,description) VALUES('parvovirus','A serious digestive disease that can cause severe diarrhea and dehydration in puppies, senior dogs, and dogs with underlying health conditions');`,
		`INSERT INTO shots(name,description) VALUES('leptospirosis','A zoonotic disease that can spread from animals to humans through contaminated water or soil');`,
		`INSERT INTO shots(name,description) VALUES('adenovirus','Canine adenovirus type 1 causes infectious canine hepatitis, a disease that affects the liver and other organs. It can lead to fever, abdominal pain, vomiting, and, in severe cases, death.');`,
		`INSERT INTO shots(name,description) VALUES('feline rhinotracheitis virus','Feline herpesvirus, also known as feline rhinotracheitis virus, causes severe signs of upper respiratory infection');`,
		`INSERT INTO shots(name,description) VALUES('feline calicivirus','FCV is thought to be associated with chronic gingivitis/stomatitis, a highly painful inflammation of the gums and teeth. ');`,
		`INSERT INTO shots(name,description) VALUES('feline panleukopenia','While the disease usually starts with decreased energy and low appetite, it progresses to vomiting and diarrhea. The virus also kills off white blood cells');`,
		`INSERT INTO shots(name,description) VALUES('FeLV','works to protect cats against the feline leukemia virus');`,
		`INSERT INTO shots(name,description) VALUES('chalmydia','causes upper respitory infections');`,
		`INSERT INTO shots(name,description) VALUES('','');`,
	}

	for _, query := range sql {
		_, err := pool.Query(ctx, query)
		queryFail(err, "populating shot data")
	}
}

func populateNotificationTypes(pool *pgxpool.Pool) {
	sql := []string{
		`INSERT INTO notification_types(name) VALUES ('request')`,
		`INSERT INTO notification_types(name) VALUES ('reply')`,
		`INSERT INTO notification_types(name) VALUES ('alert')`,
		`INSERT INTO notification_types(name) VALUES ('news')`,
	}
	for _, query := range sql {
		_, err := pool.Query(context.TODO(), query)
		queryFail(err, "populating shot data")
	}
}
