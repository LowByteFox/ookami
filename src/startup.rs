use rand::Rng;

static IMAGES: [&str; 1] = [include_str!("mio.txt")];

pub fn draw() {
    let mut rng = rand::thread_rng();
    let index = rng.gen_range(0..IMAGES.len());
    print!("{}", IMAGES[index]);
}
